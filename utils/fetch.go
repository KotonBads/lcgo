package lcgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

/*
Fetch the LunarClient launch API endpoint.

@param {string} platform - Current platform (win32, linux, darwin)

@param {string} arch - Current platform architecture (amd64, amd64p32)

@param {string} version - Minecraft version to download for (1.7, 1.8.9, ..., 1.19.4)

@returns {LaunchMeta} response - API response Unmarshalled into a struct

@returns {error} err - Any errors that occured
*/
func FetchAPI(platform string, arch string, version string, module string) (response LaunchMeta, err error) {
	url := "https://api.lunarclientprod.com/launcher/launch"

	params := map[string]string{
		"hwid":        "0",
		"os":          platform,
		"arch":        arch,
		"version":     version,
		"branch":      "master",
		"launch_type": "OFFLINE",
		"classifier":  module,
	}

	jsonVal, _ := json.Marshal(params)
	var res LaunchMeta

	if response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonVal)); err == nil {
		body, _ := io.ReadAll(response.Body)

		if err := json.Unmarshal(body, &res); err != nil {
			fmt.Printf("Couldn't Unmarshal the response: %s\n", err)
		}
		return res, err
	}
	return res, err
}
