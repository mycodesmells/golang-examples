package main

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type config struct {
	Username string `env:"USERNAME" envDefault:"Slomek"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env not found, reading configuration from ENV")
	}

	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Failed to parse ENV")
	}
	log.Printf("Hello, %s!\n", cfg.Username)
}
