package redis

import (
	"github.com/go-redis/redis/v8"
)

// RdStore is the redis store
type RdStore struct {
	client *redis.Client
}

// NewClient returns the new redis client
func NewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return client
}

// NewStore returns the new redis store
func NewStore(c *redis.Client) *RdStore {
	return &RdStore{client: c}
}
