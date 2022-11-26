package api

import (
	"context"

	"github.com/darchlabs/nodes/src/config"
	"github.com/darchlabs/nodes/src/internal/storage"
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

	for nodeID := range ctx.server.nodesCommands {
		metric := make([]*nodeMetric, 0)
		for method := range config.ETHNodesMethods {
			m, err := ctx.store.GetMethodMetric(context.Background(), &storage.GetMethodMetricInput{
				NodeID: nodeID,
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

		nodeMetrics[nodeID] = metric
	}

	return &getNodesMetricsHandlerResponse{
		Metrics: nodeMetrics,
	}, fiber.StatusOK, nil
}
