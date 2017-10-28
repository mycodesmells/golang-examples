package main

import (
	"fmt"
	"os"

	"github.com/mycodesmells/golang-examples/cobra/hello/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
