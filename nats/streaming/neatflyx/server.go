package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/protobuf/proto"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	pb "github.com/mycodesmells/golang-examples/nats/streaming/proto"
)

type publishRequest struct {
	SeriesName string `json:"series_name,omitempty"`
	SeasonNo   int    `json:"season_no,omitempty"`
	EpisodeNo  int    `json:"episode_no,omitempty"`
	EpisodeURL string `json:"episode_url,omitempty"`
}

const (
	topicPublishEpisode = "episodes:publish"
)

type server struct {
	natsClient stan.Conn
}

func (s server) HandlePublishEpisode(rw http.ResponseWriter, req *http.Request) {
	var pubReq publishRequest
	if err := json.NewDecoder(req.Body).Decode(&pubReq); err != nil {
		log.Errorf("Failed to read request: %v", err)
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	message := &pb.PublishEpisodeMessage{
		SeriesName: pubReq.SeriesName,
		SeasonNo:   int64(pubReq.SeasonNo),
		EpisodeNo:  int64(pubReq.EpisodeNo),
		EpisodeUrl: pubReq.EpisodeURL,
	}

	if err := s.publishMessage(topicPublishEpisode, message); err != nil {
		log.Errorf("Failed to publish message onto queue: %v", err)
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}

	log.Printf("Publishing on S%02dE%02d of '%s' on '%s'", message.SeasonNo, message.EpisodeNo, message.SeriesName, message.EpisodeUrl)
	fmt.Fprint(rw, "Post publication is pending")
}

func (s server) publishMessage(topic string, msg proto.Message) error {
	bs, err := proto.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal proto message")
	}

	if err := s.natsClient.Publish(topicPublishEpisode, bs); err != nil {
		return errors.Wrap(err, "failed to publish message")
	}

	return nil
}
