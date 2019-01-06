package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	listen = flag.String("listen", ":8080", "listen address")
)

func main() {
	flag.Parse()

	loop()

	log.Printf("listening on %q...", *listen)
	log.Fatal(http.ListenAndServe(*listen, http.FileServer(http.Dir("dist/"))))
}

func loop() {
	for i := 0; i < 10; i++ {
		go func(n int) {
			fmt.Printf("#%d\n", n)
		}(i)
	}
}
