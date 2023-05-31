package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	lcgo "github.com/KotonBads/lcgo/utils"
)

func main() {
	config := flag.String("config", "config.json", "Path to config file")
	debug := flag.Bool("debug", false, "Toggle debug output")
	version := flag.String("version", "", "Minecraft version to launch")

	flag.Parse()

	if _, err := os.Stat(*config); errors.Is(err, os.ErrNotExist) {
		panic(fmt.Sprintf("Config does not exist: %s\nConsider running %s -h", *config, os.Args[0]))
	}
	
	if len(*version) < 3 {
		panic(fmt.Sprintf("Specify a version\nRun %s -version <insert mc version>", os.Args[0]))
	}

	lcgo.ChangeVersion(*config, *version)

	lcgo.Launch(*config, *debug)

}
