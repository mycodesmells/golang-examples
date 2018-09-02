// The purpose of this package is to showcase how one can read password
// (or any secret) input from the command line without showing it back
// to the user.
package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	username := flag.String("u", "", "username")
	typePass := flag.Bool("p", false, "password input")
	flag.Parse()

	if *username == "" {
		fmt.Println("Username not provided")
		return
	}

	var password string
	if *typePass {
		fmt.Printf("Password for user %s: ", *username)
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Printf("Failed to read password: %v", err)
		}
		password = string(bytePassword)
		fmt.Println()
	} else {
		password = os.Getenv("MYAPP_PASSWORD")
	}

	if password == "" {
		fmt.Println("Password not provided")
		return
	}

	fmt.Printf("Credentials: %s/%s\n", *username, password)
}
