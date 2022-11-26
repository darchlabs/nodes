package storage

import (
	"context"

	"github.com/go-redis/redis/v9"
)

type DataStore interface {
	PutMethodMetric(context.Context, *PutMethodMetricInput) error
	GetMethodMetric(context.Context, *GetMethodMetricInput) (*MethodMetricRecord, error)
}

type Store struct {
	db *redis.Client
}

func NewDataStore(databaseURL string) (*Store, error) {
	db := redis.NewClient(&redis.Options{
		Addr: databaseURL,
		DB:   0, // use default DB
	})

	return &Store{
		db: db,
	}, nil
}
