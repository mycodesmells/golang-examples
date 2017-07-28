package main

import (
	"fmt"
	"os"
)

func main() {
	username := os.Getenv("USERNAME")
	if username == "" {
		username = "Slomek"
	}
	fmt.Printf("Hello, %s!\n", username)
}
