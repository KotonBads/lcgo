package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

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

func fetch(os string, arch string, version string) (natives map[string]any, err error) {
	url := "https://api.lunarclientprod.com/launcher/launch"

	params := map[string]string{
		"hwid":        "0",
		"os":          os,
		"arch":        arch,
		"version":     version,
		"branch":      "master",
		"launch_type": "OFFLINE",
		"classifier":  "optifine",
	}

	_natives := map[string]any{}

	jsonVal, _ := json.Marshal(params)

	response, error := http.Post(url, "application/json", bytes.NewBuffer(jsonVal))

	body, _ := io.ReadAll(response.Body)

	_err := json.Unmarshal(body, &_natives)

	if _err != nil {
		panic(_err)
	}

	return _natives, error
}

func write(natives map[string]any) {
	if _json, _err := json.MarshalIndent(natives, "", "\t"); _err == nil {
		if f, _err := os.Create("debug/launch.json"); _err == nil {
			f.Write(_json)
			defer f.Close()
		} else {
			log.Fatal(_err)
		}
	}
}

func main() {
	version := "1.8.9"
	os, arch := platform()

	if natives, err := fetch(os, arch, version); err == nil {
		fmt.Println(natives)
		write(natives)
	}
}
