package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/mycodesmells/golang-examples/misc/errors/hello"
)

func main() {
	if len(os.Args) == 0 {
		fmt.Println("Usage: hello <name>")
		os.Exit(1)
	}

	name := os.Args[1]
	greeting, err := hello.Hi(name)
	if err != nil {
		if errors.Is(err, hello.ErrMissingName) {
			fmt.Println("Missing name argument")
			fmt.Println("Usage: hello <name>")
			os.Exit(1)
		}

		fmt.Printf("Fatal error, so sorry but: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
}
