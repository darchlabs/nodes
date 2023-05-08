package api

import (
	"context"

	"github.com/darchlabs/nodes/config"
	"github.com/darchlabs/nodes/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type getNodesMetricsHandlerResponse struct {
	Metrics map[string][]*nodeMetric
}

type nodeMetric struct {
	Method string `json:"method"`
	Count  int64  `json:"count"`
}

func getNodesMetricsHandler(ctx *Context, c *fiber.Ctx) (interface{}, int, error) {
	nodeMetrics := make(map[string][]*nodeMetric)

	for _, instance := range ctx.server.nodesManager.GetAll() {
		metric := make([]*nodeMetric, 0)
		for method := range config.ETHNodesMethods {
			m, err := ctx.kvStore.GetMethodMetric(context.Background(), &storage.GetMethodMetricInput{
				NodeID: instance.ID,
				Method: method,
			})
			if err != nil {
				return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: getNodesMetricsHandler ctx.store.GetMethodMetric error")
			}

			metric = append(metric, &nodeMetric{
				Method: m.Method,
				Count:  m.Count,
			})
		}

		nodeMetrics[instance.ID] = metric
	}

	return &getNodesMetricsHandlerResponse{
		Metrics: nodeMetrics,
	}, fiber.StatusOK, nil
}
