package storage

import "github.com/go-redis/redis/v9"

type KVStore struct {
	db *redis.Client
}

func NewKeyValueStore(databaseURL string) (*KVStore, error) {
	db := redis.NewClient(&redis.Options{
		Addr: databaseURL,
		DB:   0, // use default DB
	})

	return &KVStore{
		db: db,
	}, nil
}
