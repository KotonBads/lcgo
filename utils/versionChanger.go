package lcgo

import (
	"encoding/json"
	"fmt"
	"os"
)

func ChangeVersion(config string, version string) {
	var ver LaunchArgs
	if file, err := os.ReadFile(config); err == nil {
		json.Unmarshal(file, &ver)

		if ver.Version != version {
			fmt.Println("Changing version to ", version)
			ver.Version = version

			if data, err := json.MarshalIndent(&ver, "", "\t"); err == nil {
				if err := os.WriteFile(config, data, 0766); err != nil {
					panic(err)
				}
				return
			}
		}
	}
}
