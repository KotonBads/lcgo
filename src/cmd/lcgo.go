package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"lcgo"
)

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

func main() {
	version := "1.8.9"
	os, arch := lcgo.Platform()

	if natives, err := fetch(os, arch, version); err == nil {
		// fmt.Println(natives)
		write(natives)
		lcgo.Launch("/home/koton-bads/Documents/Go/LC/src/cmd/config.json", true)
	}
}
