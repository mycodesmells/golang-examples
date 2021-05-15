package main

import (
	"embed"
	"fmt"
	_ "github.com/mycodesmells/golang-examples/embed/external/statik/statik" // TODO: Replace with the absolute import path
	"io"
	"os"
	"strings"

	"github.com/rakyll/statik/fs"
)

//go:embed *.txt
var assets embed.FS

var statikReadFile = func(fpath string) ([]byte, error) {
	statikFS, err := fs.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create file system: %v", err)
	}
	f, err := statikFS.Open(fpath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %v", fpath, err)
	}
	bb, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %v", fpath, err)
	}
	return bb, nil
}

var embedReadFile = func(fpath string) ([]byte, error) {
	// Statik is treating files as if they are located in OS root directory.
	return assets.ReadFile(strings.TrimLeft(fpath, "/"))
}

func main() {
	statikHelloBB, err := statikReadFile("/hello.txt")
	if err != nil {
		fmt.Printf("Failed to open file using packr: %v", err)
		os.Exit(1)
	}
	embedHelloBB, err := embedReadFile("/hello.txt")
	if err != nil {
		fmt.Printf("Failed to open file using embed: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Output from packr: %v\n", string(statikHelloBB))
	fmt.Printf("Output from embed: %v\n", string(embedHelloBB))
}
