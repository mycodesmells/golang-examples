package main

import (
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"

	pb "github.com/mycodesmells/golang-examples/nats/pubsub/proto"
)

type config struct {
	Addr       string `envconfig:"ADDR" default:":8000"`
	NatsAddr   string `envconfig:"NATS_ADDR" default:"nats://localhost:4222"`
	StaticPath string `envconfig:"STATIC_PATH" default:"./posts"`
}

const (
	topicPublishPost = "posts:publish"
)

func main() {
	// Process ENV variables
	var cfg config
	if err := envconfig.Process("bloggenerator", &cfg); err != nil {
		log.Fatalf("Failed to load configuration from env: %v", err)
	}

	// Connect to NATS
	natsClient, err := nats.Connect(cfg.NatsAddr)
	if err != nil {
		log.Fatalf("Can't connect to %s: %v\n", cfg.NatsAddr, err)
	}

	// Initialize page generator.
	gen, err := newPageGenerator(cfg.StaticPath)
	if err != nil {
		log.Fatalf("Failed to start page generator: %v", err)
	}

	// Start NATS subscriptions
	startSubscription(natsClient, topicPublishPost, generatePostPage(gen))

	// Start HTTP server serving static files
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(cfg.StaticPath)))

	log.Infof("Starting HTTP server on '%s'", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, r); err != nil {
		log.Fatal(err)
	}
}

// Start subscription and exit if failed.
func startSubscription(natsClient *nats.Conn, topic string, handler nats.MsgHandler) {
	if _, err := natsClient.Subscribe(topic, handler); err != nil {
		log.Fatalf("Failed to start subscription on '%s': %v", topic, err)
	}
	log.Infof("Started subscription on '%s'", topic)
}

// Wrapper for page generation.
func generatePostPage(gen pageGenerator) nats.MsgHandler {
	return func(natsMsg *nats.Msg) {
		log.Debug("Received new post generation queue message")

		var message pb.PublishPostMessage
		if err := proto.Unmarshal(natsMsg.Data, &message); err != nil {
			log.Errorf("Failed to unmarshal queue message: %v", err)
			return
		}

		if err := gen.Generate(message); err != nil {
			log.Errorf("Failed to generate post page: %v", err)
		}
	}
}
