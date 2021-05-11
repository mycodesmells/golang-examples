package main

import (
	"fmt"
	_ "github.com/mycodesmells/golang-examples/embed/external/statik/statik" // TODO: Replace with the absolute import path
	"io"
	"log"

	"github.com/rakyll/statik/fs"
)

func main() {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	f, err := statikFS.Open("/hello.txt")
	if err != nil {
		fmt.Printf("Failed to open file: %v", err)
		return
	}
	helloBB, err := io.ReadAll(f)
	if err != nil {
		fmt.Printf("Failed to read file: %v", err)
		return
	}

	fmt.Printf("Contents of 'hello.txt': %q\n", string(helloBB))
	fmt.Printf("Bytes of 'hello.txt': %v\n", helloBB)
}
