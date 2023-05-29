package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/cavaliergopher/grab/v3"
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

func fetch(_os string, arch string, version string) (artifacts map[string]any, err error) {
	url := "https://api.lunarclientprod.com/launcher/launch"

	params := map[string]string{
		"hwid":        "0",
		"os":          _os,
		"arch":        arch,
		"version":     version,
		"branch":      "master",
		"launch_type": "OFFLINE",
		"classifier":  "optifine",
	}

	natives := map[string]any{}

	jsonVal, _ := json.Marshal(params)

	response, error := http.Post(url, "application/json", bytes.NewBuffer(jsonVal))

	body, _ := io.ReadAll(response.Body)

	_err := json.Unmarshal(body, &natives)

	if _err != nil {
		panic(_err)
	}

	return natives, error
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

func downloadArtifacts(_os string, arch string, version string) {
	path := "offline/"

	ifExists := func(path string) bool {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			return false
		}
		return true
	}

	checkHash := func(path string, hash string) bool {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		h := sha1.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}

		if fmt.Sprintf("%x", h.Sum(nil)) == hash {
			return true
		}
		return false
	}

	url := "https://api.lunarclientprod.com/launcher/launch"

	params := map[string]string{
		"hwid":        "0",
		"os":          _os,
		"arch":        arch,
		"version":     version,
		"branch":      "master",
		"launch_type": "OFFLINE",
		"classifier":  "optifine",
	}

	type Artifact struct {
		Name string `json:"name"`
		Sha1 string `json:"sha1"`
		Url  string `json:"url"`
	}

	type LaunchMeta struct {
		Success        bool `json:"success"`
		UI             bool `json:"ui"`
		UpdateAssets   bool `json:"updateAssets"`
		LaunchTypeData struct {
			Artifacts []Artifact `json:"artifacts"`
			MainClass string     `json:"mainClass"`
		} `json:"launchTypeData"`
		Licenses []struct {
			File string `json:"file"`
			URL  string `json:"url"`
			Sha1 string `json:"sha1"`
		} `json:"licenses"`
		Textures struct {
			IndexURL  string `json:"indexUrl"`
			IndexSha1 string `json:"indexSha1"`
			BaseURL   string `json:"baseUrl"`
		} `json:"textures"`
		Jre struct {
			Download struct {
				URL       string `json:"url"`
				Extension string `json:"extension"`
			} `json:"download"`
			ExecutablePathInArchive []string    `json:"executablePathInArchive"`
			CheckFiles              [][]string  `json:"checkFiles"`
			ExtraArguments          []string    `json:"extraArguments"`
			JavawDownload           interface{} `json:"javawDownload"`
			JavawExeChecksum        interface{} `json:"javawExeChecksum"`
			JavaExeChecksum         string      `json:"javaExeChecksum"`
		} `json:"jre"`
	}

	jsonVal, _ := json.Marshal(params)
	var natives LaunchMeta

	response, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonVal))

	body, _ := io.ReadAll(response.Body)
	_err := json.Unmarshal(body, &natives)

	if _err != nil {
		panic(_err)
	}

	for _, v := range natives.LaunchTypeData.Artifacts {
		if !ifExists(fmt.Sprintf("%s/%s", path, v.Name)) || !checkHash(fmt.Sprintf("%s/%s", path, v.Name), v.Sha1) {
			file, err := grab.Get(fmt.Sprintf("%s/%s", path, v.Name), v.Url)
			if err != nil {
				panic(err)
			}
			fmt.Println("Downloaded file: ", file.Filename)
		} else {
			fmt.Println(v.Name, ": Already Downloaded / Up to date")
		}
	}
}

func main() {
	version := "1.8.9"
	os, arch := platform()

	if natives, err := fetch(os, arch, version); err == nil {
		// fmt.Println(natives)
		write(natives)
		downloadArtifacts(os, arch, version)
	}
}
