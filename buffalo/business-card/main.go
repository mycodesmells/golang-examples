package main

import (
	"log"

	"github.com/mycodesmells/golang-examples/buffalo/business-card/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
