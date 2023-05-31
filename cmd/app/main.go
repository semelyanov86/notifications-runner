package main

import (
	"Notifications/internal/app"
	"flag"
	"fmt"
	"log"
	"os"
)

var version string
var buildTime string

func main() {
	dirname, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	log.Println(dirname)
	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}
	app.Run(dirname)
}
