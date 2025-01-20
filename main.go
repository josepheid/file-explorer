package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/goteleport-interview/fs4/api"
)

const listenPort = 8080

//go:embed web/dist
var assets embed.FS

func main() {
	webassets, err := fs.Sub(assets, "web/dist")
	if err != nil {
		log.Fatalln("could not embed webassets", err)
	}

	rootPath := os.Getenv("ROOT_PATH")
	if rootPath == "" {
		rootPath = "./"
	}

	// check if root directory exists and is a directory
	root, err := os.Stat(rootPath)
	if err != nil {
		log.Fatalf("Root directory not found: %s\n", err)
	}
	if !root.IsDir() {
		log.Fatal("Root directory not a directory: ", rootPath)
	}

	s, err := api.NewServer(webassets, rootPath)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("starting server on port", listenPort)
	log.Fatalln(s.ListenAndServe(fmt.Sprintf("localhost:%d", listenPort)))
}
