package lcgo

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
	"path/filepath"

	"github.com/xyproto/unzip"

	"github.com/cavaliergopher/grab/v3"
)

type Artifact struct {
	Name string `json:"name"`
	Sha1 string `json:"sha1"`
	Url  string `json:"url"`
}

type LaunchMeta struct {
	Success        bool `json:"success"`
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

/*
Download artifacts from LunarClient's API.

@param {string} platform - Current platform (win32, linux, darwin)

@param {string} arch - Current platform architecture (amd64, amd64p32)

@param {string} version - Minecraft version to download for (1.7, 1.8.9, ..., 1.19.4)

@param {string} path - Download path (~/.lunarclient/offline)
*/
func downloadArtifacts(platform string, arch string, version string, path string) {

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
		"os":          platform,
		"arch":        arch,
		"version":     version,
		"branch":      "master",
		"launch_type": "OFFLINE",
		"classifier":  "optifine",
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
		if !ifExists(fmt.Sprintf("%s/%s", path, v.Name)) || !checkHash(fmt.Sprintf("%s/%s", path, v.Name), v.Sha1) || filepath.Ext(fmt.Sprintf("%s/%s", path, v.Name)) == ".zip" {
			file, err := grab.Get(fmt.Sprintf("%s/%s", path, v.Name), v.Url)
			if err != nil {
				panic(err)
			}

			fmt.Println("Downloaded file: ", file.Filename)

			if filepath.Ext(fmt.Sprintf("%s/%s", path, v.Name)) == ".zip" {
				fmt.Println("Found a zip file, Extracting...")
				if err := unzip.Extract(fmt.Sprintf("%s/%s", path, v.Name), fmt.Sprintf("%s/natives", path)); err != nil {
					panic(err)
				}
			}
		} else {
			fmt.Println(v.Name, ": Already Downloaded / Up to date")
		}
	}
}
