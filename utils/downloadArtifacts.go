package lcgo

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/xyproto/unzip"

	"github.com/cavaliergopher/grab/v3"
)

/*
Download artifacts from LunarClient's API.

@param {string} platform - Current platform (win32, linux, darwin)

@param {string} arch - Current platform architecture (amd64, amd64p32)

@param {string} version - Minecraft version to download for (1.7, 1.8.9, ..., 1.19.4)

@param {string} path - Download path (~/.lunarclient/offline)

@returns {[]Artifacts} - Array of Artifacts
*/
func DownloadArtifacts(platform string, arch string, version string, path string, module string) (artifacts []Artifacts) {
	_file := "/"

	if platform == "win32" {
		_file = "\\"
	}

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

	natives, err := FetchAPI(platform, arch, version, module)

	if err != nil {
		panic(err)
	}

	for _, v := range natives.LaunchTypeData.Artifacts {
		if !ifExists(fmt.Sprintf("%s%s%s", path, _file, v.Name)) || !checkHash(fmt.Sprintf("%s%s%s", path, _file, v.Name), v.Sha1) || filepath.Ext(fmt.Sprintf("%s%s%s", path, _file, v.Name)) == ".zip" {
			file, err := grab.Get(fmt.Sprintf("%s%s%s", path, _file, v.Name), v.Url)
			if err != nil {
				fmt.Printf("Couldn't download %s: %s\n", v.Url, err)
			}

			fmt.Println("Downloaded file: ", file.Filename)

			if filepath.Ext(fmt.Sprintf("%s%s%s", path, _file, v.Name)) == ".zip" {
				fmt.Println("Found a zip file, Extracting...")
				if err := unzip.Extract(fmt.Sprintf("%s%s%s", path, _file, v.Name), fmt.Sprintf("%s%snatives", path, _file)); err != nil {
					panic(err)
				}
			}
		} else {
			fmt.Println(v.Name, ": Already Downloaded / Up to date")
		}
	}
	return natives.LaunchTypeData.Artifacts
}
