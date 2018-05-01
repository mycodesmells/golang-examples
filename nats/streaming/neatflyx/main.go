package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	stan "github.com/nats-io/go-nats-streaming"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var (
	clusterID string = "test-cluster"
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

	// Connect to NATS-Streaming
	natsClient, err := stan.Connect(clusterID, uuid.NewV4().String(), stan.NatsURL(cfg.NatsAddr))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, cfg.NatsAddr)
	}
	defer natsClient.Close()

	srv := server{
		natsClient: natsClient,
	}

	// Serve HTTP
	r := mux.NewRouter()
	r.HandleFunc("/publish", srv.HandlePublishEpisode)

	log.Infof("Starting HTTP server on '%s'", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, r); err != nil {
		log.Fatal(err)
	}
}
