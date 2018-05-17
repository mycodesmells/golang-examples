package main

import (
	"database/sql"
	"flag"
	"log"

	proto "github.com/gogo/protobuf/proto"
	_ "github.com/lib/pq"
	mypb "github.com/mycodesmells/golang-examples/nats/streaming/proto"
	"github.com/nats-io/go-nats-streaming/pb"
	"github.com/nats-io/nats-streaming-server/spb"
	"github.com/pkg/errors"
)

var queries = map[string]string{
	"message":      "SELECT data FROM messages WHERE id=$1",
	"subscription": "SELECT proto FROM subscriptions WHERE id=$1",
	"serverinfo":   "SELECT proto FROM serverinfo WHERE id=$1 OR id='test-cluster'",
}

type protoParser func([]byte) error

var parsers = map[string]protoParser{
	"message":      parseMessage,
	"subscription": parseSubscriptions,
	"serverinfo":   parseServerInfo,
}

var (
	objectType = flag.String("obj-type", "", "NATS object type to decode from the database")
	entityID   = flag.Int("id", 0, "id of the entity to decode")
)

func main() {
	flag.Parse()

	db, err := sql.Open("postgres", "host=localhost port=15432 user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	query, ok := queries[*objectType]
	if !ok {
		log.Fatalf("No SQL query defined for object type '%s'", *objectType)
	}
	parser, ok := parsers[*objectType]
	if !ok {
		log.Fatalf("No proto parser defined for object type '%s'", *objectType)
	}

	processRows(db, *entityID, query, parser)
}

func processRows(db *sql.DB, id int, query string, parser protoParser) {
	rows, err := db.Query(query, id)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pdata []byte
		if err := rows.Scan(&pdata); err != nil {
			log.Printf("Failed to read binary column value: %v\n", err)
			continue
		}

		if err := parser(pdata); err != nil {
			log.Printf("Failed to process proto data: %v\n", err)
			continue
		}
	}
}

func parseMessage(data []byte) error {
	var p pb.MsgProto
	if err := proto.Unmarshal(data, &p); err != nil {
		return errors.Wrap(err, "failed to unmarshal proto")
	}

	var msg mypb.PublishEpisodeMessage
	if err := proto.Unmarshal(p.Data, &msg); err != nil {
		return errors.Wrap(err, "failed to unmarshal message data")
	}

	log.Println("Read message:")
	log.Printf("Subject: %s\n", p.Subject)
	log.Printf("Timestamp: %d\n", p.Timestamp)
	log.Printf("Message: %+v\n", msg)
	log.Println("=============")

	return nil
}

func parseSubscriptions(data []byte) error {
	var sub pb.SubscriptionResponse
	if err := proto.Unmarshal(data, &sub); err != nil {
		return errors.Wrap(err, "failed to unmarshal proto")
	}

	log.Println("Read subscription:")
	log.Printf("Ack Inbox: %+v\n", sub.AckInbox)
	log.Println("=============")

	return nil
}

func parseServerInfo(data []byte) error {
	// var sub pb.Ser
	var si spb.ServerInfo
	if err := proto.Unmarshal(data, &si); err != nil {
		return errors.Wrap(err, "failed to unmarshal proto")
	}

	log.Println("Read server info:")
	log.Printf("Publish: %+v\n", si.Publish)
	log.Printf("Cluster ID: %+v\n", si.ClusterID)
	log.Printf("Discovery: %+v\n", si.Discovery)
	log.Printf("Subscribe: %+v\n", si.Subscribe)
	log.Println("=============")

	return nil
}
