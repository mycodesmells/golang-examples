package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
)

type config struct {
	Addr     string `envconfig:"ADDR" default:":8000"`
	NatsAddr string `envconfig:"NATS_ADDR" default:"nats://localhost:4222"`
}

func main() {
	// Process ENV variables
	var cfg config
	if err := envconfig.Process("blogadmin", &cfg); err != nil {
		log.Fatalf("Failed to load configuration from env: %v", err)
	}

	// Connect to NATS
	natsClient, err := nats.Connect(cfg.NatsAddr)
	if err != nil {
		log.Fatalf("Can't connect to %s: %v\n", cfg.NatsAddr, err)
	}

	srv := server{
		natsClient: natsClient,
	}

	// Serve HTTP
	r := mux.NewRouter()
	r.HandleFunc("/publish", srv.HandlePublishPost)

	log.Infof("Starting HTTP server on '%s'", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, r); err != nil {
		log.Fatal(err)
	}
}
