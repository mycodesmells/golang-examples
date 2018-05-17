package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/protobuf/proto"
	"github.com/kelseyhightower/envconfig"
	stan "github.com/nats-io/go-nats-streaming"
	stanpb "github.com/nats-io/go-nats-streaming/pb"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	pb "github.com/mycodesmells/golang-examples/nats/streaming/proto"
)

type config struct {
	NatsAddr string `envconfig:"NATS_ADDR" default:"nats://localhost:4222"`
	StartOpt string `envconfig:"START_OPT" default:"ONLY_NEW"`
}

var (
	clusterID           = "test-cluster"
	topicPublishEpisode = "episodes:publish"
)

func main() {
	// Process ENV variables
	var cfg config
	if err := envconfig.Process("bloggenerator", &cfg); err != nil {
		log.Fatalf("Failed to load configuration from env: %v", err)
	}

	// Connect to NATS
	natsClient, err := stan.Connect(clusterID, uuid.NewV4().String(), stan.NatsURL(cfg.NatsAddr))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, cfg.NatsAddr)
	}
	defer natsClient.Close()

	// Start NATS subscriptions
	startSubscription(natsClient, topicPublishEpisode, watchEpisode, startOpt(cfg.StartOpt))

	log.Infof("Starting new watcher service")

	// Waiting for signal to shutdown.
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func startSubscription(natsClient stan.Conn, topic string, handler stan.MsgHandler, startOpt stan.SubscriptionOption) {
	durableName := uuid.NewV4().String()

	if _, err := natsClient.QueueSubscribe(topic, durableName, handler, startOpt, stan.DurableName(durableName)); err != nil {
		natsClient.Close()
		log.Fatal(err)
	}
	log.Infof("Started new regular subscription")
}

func watchEpisode(natsMsg *stan.Msg) {
	log.Debug("Received new post generation queue message")

	var message pb.PublishEpisodeMessage
	if err := proto.Unmarshal(natsMsg.Data, &message); err != nil {
		log.Errorf("Failed to unmarshal queue message: %v", err)
		return
	}

	log.Printf("Watching on S%02dE%02d of '%s' on '%s'", message.SeasonNo, message.EpisodeNo, message.SeriesName, message.EpisodeUrl)
}

func startOpt(optString string) stan.SubscriptionOption {
	switch optString {
	default:
		return stan.StartAt(stanpb.StartPosition_NewOnly)
	case "MOST_RECENT":
		return stan.StartWithLastReceived()
	case "ALL":
		return stan.DeliverAllAvailable()
	}
}
