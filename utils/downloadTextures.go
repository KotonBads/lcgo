package lcgo

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/cavaliergopher/grab/v3"
)

func DownloadTextures(platform string, arch string, version string, path string, debug bool, module string) {
	var response LaunchMeta
	var wg sync.WaitGroup

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

	download := func(index []string) {
		for _, v := range index {
			wg.Add(1)
			go func(v string) {
				defer wg.Done()
				e := strings.Split(v, " ")

				if !ifExists(path+_file+e[0]) || !checkHash(path+_file+e[0], e[1]) {
					file, err := grab.Get(path+_file+e[0], response.Textures.BaseURL+e[1])

					if err != nil {
						panic(err)
					}

					if debug {
						fmt.Println()
						fmt.Println("Downloaded: ", file.Filename)
					}
				} else {
					if debug {
						fmt.Println()
						fmt.Println("Already up-to-date: ", e[0])
					}
				}
			}(v)
		}
		wg.Wait()
	}

	if res, err := FetchAPI(platform, arch, version, module); err == nil {
		response = res
	}

	if res, err := http.Get(response.Textures.IndexURL); err == nil {
		if body, err := io.ReadAll(res.Body); err == nil {
			download(strings.Split(string(body), "\n"))
		}
	}
}
