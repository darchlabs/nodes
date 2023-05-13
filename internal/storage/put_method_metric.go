package storage

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
)

type PutMethodMetricInput struct {
	NodeID string
	Method string
}

func (s *KVStore) PutMethodMetric(ctx context.Context, input *PutMethodMetricInput) error {
	key := fmt.Sprintf("%s-%s", input.NodeID, input.Method)

	result, err := s.db.Get(ctx, key).Result()
	switch true {
	case errors.Is(err, redis.Nil):
		result = "0"
	case err != nil:
		return errors.Wrap(err, "storage: Store.PutMethodMetric s.db.Get error")
	}

	value, err := strconv.Atoi(result)
	if err != nil {
		return errors.Wrap(err, "storage: Store.PutMethodMetric strconv.Atoi error")
	}

	value++
	err = s.db.Set(ctx, key, value, 0).Err()
	if err != nil {
		return errors.Wrap(err, "storage: Store.PutMethodMetric s.db.Set error")
	}

	return nil
}
