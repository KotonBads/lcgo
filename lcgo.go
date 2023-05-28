package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func fetch() (natives map[string]any, err error) {
	url := "https://api.lunarclientprod.com/launcher/launch"

	params := map[string]string{
		"hwid":        "0",
		"os":          "win32",
		"arch":        "x64",
		"version":     "1.8.9",
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
	if natives, err := fetch(); err == nil {
		fmt.Println(natives)
		write(natives)
	}
}
