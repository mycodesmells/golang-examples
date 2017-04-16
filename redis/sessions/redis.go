package sessions

import (
	"encoding/json"
	"log"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type redisStore struct {
	client *redis.Client
}

func NewRedisStore() Store {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}

	return &redisStore{
		client: client,
	}
}

func (r redisStore) Get(id string) (Session, error) {
	var session Session

	bs, err := r.client.Get(id).Bytes()
	if err != nil {
		return session, errors.Wrap(err, "failed to get session from redis")
	}

	if err := json.Unmarshal(bs, &session); err != nil {
		return session, errors.Wrap(err, "failed to unmarshall session data")
	}

	return session, nil
}

func (r redisStore) Set(id string, session Session) error {
	bs, err := json.Marshal(session)
	if err != nil {
		return errors.Wrap(err, "failed to save session to redis")
	}

	if err := r.client.Set(id, bs, 0).Err(); err != nil {
		return errors.Wrap(err, "failed to save session to redis")
	}

	return nil
}
