package main

import (
	"log"

	"github.com/gobuffalo/envy"
	"github.com/mycodesmells/golang-examples/buffalo/business-card/actions"
)

func main() {
	port := envy.Get("PORT", "3000")
	app := actions.App()
	log.Fatal(app.Start(port))
}
