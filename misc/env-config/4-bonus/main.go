package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("printenv", "USERNAME")
	cmd.Env = append(os.Environ(), "USERNAME=Oscar")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
