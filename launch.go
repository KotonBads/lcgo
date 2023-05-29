package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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
	Assets   string   `json:"assets"`
	Textures string   `json:"textures"`
	Natives  string   `json:"natives"`
	MCDir    string   `json:"mcdir"`
	Agents   []string `json:"agents"`
	JVMArgs  []string `json:"JVM Args"`
}

func launch(config string) {
	var launchArgs LaunchArgs
	var assets []string
	var assetsPath []string
	var assetIndex string

	ichorClassPath := func() string {
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
		return fmt.Sprintf("%s/%s,%s", launchArgs.Assets, e, ver)
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

	if launchArgs.Version == "1.8.9" {
		assetIndex = "1.8"
	} else {
		assetIndex = launchArgs.Version
	}

	if f, err := os.ReadFile(config); err == nil {
		if err != nil {
			panic(err)
		}

		json.Unmarshal(f, &launchArgs)
	}

	if entries, err := os.ReadDir(launchArgs.Assets); err == nil {
		for _, v := range entries {
			if !v.IsDir() && !strings.HasSuffix(v.Name(), ".zip") {
				assets = append(assets, v.Name())
				assetsPath = append(assetsPath, fmt.Sprintf("\"%s/%s\"", launchArgs.Assets, v.Name()))
			}
		}
	}

	out, err := exec.Command(fmt.Sprintf("%s/bin/java --add-modules jdk.naming.dns --add-exports jdk.naming.dns/com.sun.jndi.dns=java.naming -Djava.library.path=%s -Dlog4j2.formatMsgNoLookups=true --add-opens java.base/java.io=ALL-UNNAMED -Xms%s -Xmx%s -Xss%s -Xmn%s %s -cp %s %s com.moonsworth.lunar.genesis.Genesis --version %s --accessToken 0 --assetIndex %s --userProperties {} --gameDir %s --texturesDir %s --launcherVersion 69420 --hwid 69420 --width %d --height %d --workingDirectory %s --classpathDir %s --ichorClassPath \"%s\"", launchArgs.JRE, launchArgs.Natives, launchArgs.Memory.Xms, launchArgs.Memory.Xmx, launchArgs.Memory.Xss, launchArgs.Memory.Xmn, strings.Join(launchArgs.JVMArgs, " "), strings.Join(assetsPath, ":"), javaAgent(), launchArgs.Version, assetIndex, launchArgs.MCDir, launchArgs.Textures, launchArgs.Width, launchArgs.Height, launchArgs.Assets, launchArgs.Assets, ichorClassPath())).CombinedOutput()

	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}
