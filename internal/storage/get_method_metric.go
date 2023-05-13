package storage

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
)

type MethodMetricRecord struct {
	NodeID string
	Method string
	Count  int64
}

type GetMethodMetricInput struct {
	Method string
	NodeID string
}

func (s *KVStore) GetMethodMetric(ctx context.Context, input *GetMethodMetricInput) (*MethodMetricRecord, error) {
	key := fmt.Sprintf("%s-%s", input.NodeID, input.Method)
	result, err := s.db.Get(ctx, key).Result()
	switch true {
	case errors.Is(err, redis.Nil):
		result = "0"
	case err != nil:
		return nil, errors.Wrap(err, "storage: Store.GetMethodMetrics s.db.Get error")
	}

	value, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "storage: Store.GetMethodMetric strconv.Atoi error")
	}

	return &MethodMetricRecord{
		NodeID: input.NodeID,
		Method: input.Method,
		Count:  value,
	}, nil
}
