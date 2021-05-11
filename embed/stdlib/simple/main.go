package main

import (
	_ "embed"
	"fmt"
)

//go:embed hello.txt
var hello string

//go:embed hello.txt
var helloBB []byte

func main() {
	fmt.Printf("Contents of 'hello.txt': %q\n", hello)
	fmt.Printf("Bytes of 'hello.txt': %v\n", helloBB)
}
