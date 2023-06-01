package lcgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

/*
Returns the current Platform and architecture. Panics when LunarClient doesn't support the Platform and/or the architecture.

@returns {string} Platform - Current Platform (win32, linux, darwin)

@returns {string} arch - Current Platform architecture (amd64, amd64p32)
*/
func Platform() (platform string, arch string) {
	var _platform string
	var _arch string

	switch runtime.GOOS {
	case "windows":
		_platform = "win32"
	case "linux":
		_platform = "linux"
	case "darwin":
		_platform = "darwin"
	default:
		panic(fmt.Sprintf("Platform not supported: %s", runtime.GOOS))
	}

	switch runtime.GOARCH {
	case "amd64":
		_arch = "x64"
	case "amd64p32":
		if _platform != "win32" {
			panic(fmt.Sprintf("Arch not supported for platform: %s", _platform))
		}
		_arch = "ia32"
	default:
		panic(fmt.Sprintf("Arch not supported: %s", runtime.GOARCH))
	}
	return _platform, _arch
}

/*
Launches LunarClient.

@param {string} config - Path to config file

@param {bool} debug - Toggles debug output
*/
func Launch(config string, debug bool) {
	var launchArgs LaunchArgs
	var classPath []string
	var externalFiles []string
	var assetIndex string
	sep := ":"

	javaAgent := func() string {
		var e string

		if len(launchArgs.Agents) > 0 {
			for _, v := range launchArgs.Agents {
				if len(launchArgs.Agents) > 1 {
					e += fmt.Sprintf("-javaagent:\"%s\" ", v)
				} else {
					e += fmt.Sprintf("-javaagent:\"%s\"", v)
				}
			}
		}
		return e
	}

	fallbackPath := func(path string) (fallback string) {
		fmt.Println() // spacing
		folder, err := os.Stat(path)
		if path == "" {
			fmt.Printf("Empty Path: %s\n", err)
			fmt.Printf("Falling back to downloadPath: %s\n\n", launchArgs.DownloadPath)
			return launchArgs.DownloadPath + "/" + launchArgs.Version
		}
		if !folder.IsDir() {
			if folder.Name() != path {
				fmt.Printf("Path is not a folder: %s\n", err)
				fmt.Printf("Falling back to downloadPath: %s\n", launchArgs.DownloadPath)
				return launchArgs.DownloadPath + "/" + launchArgs.Version
			}
			fmt.Printf("Path is not a folder: %s\n", err)
			panic("downloadPath and path is the same! No folder to fall back to.")
		}
		if os.IsNotExist(err) {
			if folder.Name() != path {
				fmt.Printf("Folder does not exist: %s\n", err)
				fmt.Printf("Falling back to downloadPath: %s\n", launchArgs.DownloadPath)
				return launchArgs.DownloadPath + "/" + launchArgs.Version
			}
			fmt.Printf("Folder does not exist: %s\n", err)
			panic("downloadPath and path is the same! No folder to fall back to.")
		}
		return launchArgs.DownloadPath + "/" + launchArgs.Version
	}

	for _, v := range launchArgs.Env {
		if err := os.Setenv(v.Name, v.Value); err != nil {
			panic(err)
		}
	}

	if f, err := os.ReadFile(config); err == nil {
		if err != nil {
			panic(err)
		}

		json.Unmarshal(f, &launchArgs)
	}

	if launchArgs.Version == "1.8.9" {
		assetIndex = "1.8"
	} else {
		assetIndex = launchArgs.Version
	}

	if len(launchArgs.PreJava) > 0 {
		launchArgs.JRE = fmt.Sprintf("%s %s", launchArgs.PreJava, launchArgs.JRE)
	}

	launchArgs.Natives = fmt.Sprintf("\"%s/natives\"", fallbackPath(launchArgs.Natives))
	launchArgs.Assets = fallbackPath(launchArgs.Assets)

	plat, arch := Platform()

	artifacts := downloadArtifacts(plat, arch, launchArgs.Version, launchArgs.Assets)

	for _, v := range artifacts {
		if v.Type == "CLASS_PATH" {
			classPath = append(classPath, fmt.Sprintf("%s/%s", launchArgs.Assets, v.Name))
		}

		if v.Type == "EXTERNAL_FILE" {
			externalFiles = append(externalFiles, fmt.Sprintf("%s/%s", launchArgs.Assets, v.Name))
		}
	}

	if plat == "win32" {
		sep = ";"
	}

	cmd := exec.Command("bash", "-c", fmt.Sprintf("%s/bin/java --add-modules jdk.naming.dns --add-exports jdk.naming.dns/com.sun.jndi.dns=java.naming -Djna.boot.library.path=%s -Djava.library.path=%s -Dlog4j2.formatMsgNoLookups=true --add-opens java.base/java.io=ALL-UNNAMED -Xms%s -Xmx%s -Xss%s -Xmn%s %s -cp %s %s com.moonsworth.lunar.genesis.Genesis --version %s --accessToken 0 --assetIndex %s --userProperties {} --gameDir %s --texturesDir %s --launcherVersion 69420 --hwid 69420 --width %d --height %d --workingDirectory %s --classpathDir %s --ichorClassPath %s --ichorExternalFiles %s", launchArgs.JRE, launchArgs.Natives, launchArgs.Natives, launchArgs.Memory.Xms, launchArgs.Memory.Xmx, launchArgs.Memory.Xss, launchArgs.Memory.Xmn, strings.Join(launchArgs.JVMArgs, " "), strings.Join(classPath, sep), javaAgent(), launchArgs.Version, assetIndex, launchArgs.MCDir, launchArgs.Textures, launchArgs.Width, launchArgs.Height, launchArgs.Assets, launchArgs.Assets, strings.Join(classPath, ","), strings.Join(externalFiles, ",")))

	if debug {
		fmt.Println()
		fmt.Printf("Config: %s\n", config)
		fmt.Printf("Platform: %s %s\n", plat, arch)
		fmt.Printf("MC Version: %s\nAsset Index: %s\nIchor: %s\n", launchArgs.Version, assetIndex, strings.Split(cmd.Args[2], ",")[1])
		fmt.Printf("Using JRE: %s\n", launchArgs.JRE)
		fmt.Printf("Natives: %s\n", launchArgs.Natives)
		fmt.Printf("Assets: %s\n", launchArgs.Assets)
		fmt.Printf("PreJava: %s\n", launchArgs.PreJava)
		fmt.Printf("Env: %s\n", launchArgs.Env)

		fmt.Printf("\nExecuting: \n%s\n\n", strings.Join(cmd.Args, " "))
	}

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
