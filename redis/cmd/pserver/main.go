package main

import (
	"flag"
	"log"

	redis "gopkg.in/redis.v4"
)

var task string

func init() {
	flag.StringVar(&task, "task", "", "task to be handled by client")
}

func main() {
	flag.Parse()

	if task == "" {
		log.Fatal("Task must not be empty")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}

	if err := client.Publish("pubsub-key", task).Err(); err != nil {
		log.Fatalf("Failed to put stuff into queue: %v", err)
	}
	log.Printf("'%v' task put into queue", task)
}
