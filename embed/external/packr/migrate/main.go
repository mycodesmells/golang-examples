package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr"
)

//go:embed *.txt
var assets embed.FS

var packrReadFile = func(fpath string) ([]byte, error) {
	fname := filepath.Base(fpath)
	box := packr.NewBox(strings.TrimRight(fpath, fname))
	bb, err := box.Find(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", fpath, err)
	}
	return bb, nil
}

var embedReadFile = assets.ReadFile

func main() {
	packrHelloBB, err := packrReadFile("hello.txt")
	if err != nil {
		fmt.Printf("Failed to open file using packr: %v", err)
		os.Exit(1)
	}
	embedHelloBB, err := embedReadFile("hello.txt")
	if err != nil {
		fmt.Printf("Failed to open file using embed: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Output from packr: %v\n", string(packrHelloBB))
	fmt.Printf("Output from embed: %v\n", string(embedHelloBB))
}
