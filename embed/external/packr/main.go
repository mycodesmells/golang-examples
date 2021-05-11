package main

import (
	"fmt"

	"github.com/gobuffalo/packr"
)

func main() {
	box := packr.NewBox(".")
	hello, err := box.FindString("hello.txt")
	if err != nil {
		fmt.Printf("Failed to find string file: %v", err)
		return
	}

	fmt.Printf("Contents of 'hello.txt': %q\n", hello)
	fmt.Printf("Bytes of 'hello.txt': %v\n", []byte(hello))
}
