package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
)

type config struct {
	Username string `env:"USERNAME" envDefault:"Slomek"`
}

func main() {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse ENV")
	}
	fmt.Printf("Hello, %s!\n", cfg.Username)
}
