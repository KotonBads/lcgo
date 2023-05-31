package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/KotonBads/lcgo/utils"
	"os"
)

func main() {
	config := flag.String("config", "config.json", "Path to config file")
	debug := flag.Bool("debug", false, "Toggle debug output")

	flag.Parse()

	if _, err := os.Stat(*config); errors.Is(err, os.ErrNotExist) {
		panic(fmt.Sprintf("Config does not exist: %s\nConsider running %s -h", *config, os.Args[0]))
	}

	lcgo.Launch(*config, *debug)

}
