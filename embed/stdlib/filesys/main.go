package main

import (
	"embed"
	_ "embed"
	"fmt"
)

//go:embed assets/*.txt
var assets embed.FS

func main() {
	file, err := assets.ReadFile("assets/bye.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Contents of 'assets/bye.txt': %q\n", file)
}
