package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/protobuf/proto"
	pb "github.com/mycodesmells/golang-examples/nats/proto"
	nats "github.com/nats-io/go-nats"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type publishRequest struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

const (
	topicPublishPost = "posts:publish"
)

type server struct {
	natsClient *nats.Conn
}

func (s server) HandlePublishPost(rw http.ResponseWriter, req *http.Request) {
	var pubReq publishRequest
	if err := json.NewDecoder(req.Body).Decode(&pubReq); err != nil {
		log.Errorf("Failed to read request: %v", err)
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	message := &pb.PublishPostMessage{
		Title:   pubReq.Title,
		Content: pubReq.Content,
	}

	if err := s.publishMessage(topicPublishPost, message); err != nil {
		log.Errorf("Failed to publish message onto queue: %v", err)
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}

	log.Printf("Publishing on '%s': '%s'", topicPublishPost, pubReq.Content)
	fmt.Fprint(rw, "Post publication is pending")
}

func (s server) publishMessage(topic string, msg proto.Message) error {
	bs, err := proto.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal proto message")
	}

	if err := s.natsClient.Publish(topicPublishPost, bs); err != nil {
		return errors.Wrap(err, "failed to publish message")
	}

	if err := s.natsClient.Flush(); err != nil {
		return errors.Wrap(err, "failed to flush message")
	}

	if err := s.natsClient.LastError(); err != nil {
		return errors.Wrap(err, "received error after publishing")
	}

	return nil
}
