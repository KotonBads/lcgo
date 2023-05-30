package main

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

type LaunchArgs struct {
	Version string `json:"version"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	JRE     string `json:"jre"`
	Memory  struct {
		Xms string `json:"Xms"`
		Xmx string `json:"Xmx"`
		Xss string `json:"Xss"`
		Xmn string `json:"Xmn"`
	}
	DownloadPath string   `json:"downloadPath"`
	Assets       string   `json:"assets"`
	Textures     string   `json:"textures"`
	Natives      string   `json:"natives"`
	MCDir        string   `json:"mcdir"`
	Agents       []string `json:"agents"`
	JVMArgs      []string `json:"JVM Args"`
}

/*
Returns the current platform and architecture. Panics when LunarClient doesn't support the platform and/or the architecture.

@returns {string} platform - Current platform (win32, linux, darwin)

@returns {string} arch - Current platform architecture (amd64, amd64p32)
*/
func platform() (platform string, arch string) {
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
func launch(config string, debug bool) {
	var launchArgs LaunchArgs
	var assets []string
	var assetsPath []string
	assetIndex := launchArgs.Version

	ichorClassPath := func(path string) string {
		var e string
		var ver string

		for _, v := range assets {
			if strings.HasPrefix(v, "genesis") {
				e += v
			}

			if strings.HasPrefix(v, "v1_") {
				ver += strings.Split(v, "-")[0]
			}
		}
		return fmt.Sprintf("%s/%s,%s", path, e, ver)
	}

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

	if f, err := os.ReadFile(config); err == nil {
		if err != nil {
			panic(err)
		}

		json.Unmarshal(f, &launchArgs)
	}

	if launchArgs.Version == "1.8.9" {
		assetIndex = "1.8"
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

	launchArgs.Natives = fmt.Sprintf("\"%s/natives\"", fallbackPath(launchArgs.Natives))
	launchArgs.Assets = fallbackPath(launchArgs.Assets)

	plat, arch := platform()

	downloadArtifacts(plat, arch, launchArgs.Version, launchArgs.Assets)

	if entries, err := os.ReadDir(launchArgs.Assets); err == nil {
		for _, v := range entries {
			if !v.IsDir() && !strings.HasSuffix(v.Name(), ".zip") && !strings.HasPrefix(v.Name(), "OptiFine") {
				assets = append(assets, v.Name())
				assetsPath = append(assetsPath, fmt.Sprintf("%s/%s", launchArgs.Assets, v.Name()))
			}
		}
	}

	cmd := exec.Command("bash", "-c", fmt.Sprintf("%s/bin/java --add-modules jdk.naming.dns --add-exports jdk.naming.dns/com.sun.jndi.dns=java.naming -Djna.boot.library.path=%s -Djava.library.path=%s -Dlog4j2.formatMsgNoLookups=true --add-opens java.base/java.io=ALL-UNNAMED -Xms%s -Xmx%s -Xss%s -Xmn%s %s -cp %s %s com.moonsworth.lunar.genesis.Genesis --version %s --accessToken 0 --assetIndex %s --userProperties {} --gameDir %s --texturesDir %s --launcherVersion 69420 --hwid 69420 --width %d --height %d --workingDirectory %s --classpathDir %s --ichorClassPath %s", launchArgs.JRE, launchArgs.Natives, launchArgs.Natives, launchArgs.Memory.Xms, launchArgs.Memory.Xmx, launchArgs.Memory.Xss, launchArgs.Memory.Xmn, strings.Join(launchArgs.JVMArgs, " "), strings.Join(assetsPath, ":"), javaAgent(), launchArgs.Version, assetIndex, launchArgs.MCDir, launchArgs.Textures, launchArgs.Width, launchArgs.Height, launchArgs.Assets, launchArgs.Assets, ichorClassPath(launchArgs.Assets)))

	if debug {
		fmt.Printf("Platform: %s %s\n", plat, arch)
		fmt.Printf("Versions: %s %s %s\n", launchArgs.Version, assetIndex, strings.Split(cmd.Args[2], ",")[1])
		fmt.Printf("Using JRE: %s\n", launchArgs.JRE)
		fmt.Printf("Natives: %s\n", launchArgs.Natives)
		fmt.Printf("Assets: %s\n", launchArgs.Assets)

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
