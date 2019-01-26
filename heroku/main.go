package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, "This is mycodesmells/golang-examples server from Heroku!")
	})
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, nil)
}
