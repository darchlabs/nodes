package api

import (
	"context"
	"fmt"
	"log"

	"github.com/darchlabs/nodes/internal/manager"
	"github.com/darchlabs/nodes/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/pkg/errors"
)

var statusAlreadyProxied int = 1000

type proxyRpcHandlerRequest struct {
	JSONRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      interface{}   `json:"id"`
}

func proxyFunc(ctx *Context) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		nodeID := c.Params("node_id")

		nodeInstance, err := ctx.server.nodesManager.Get(nodeID)
		if errors.Is(err, manager.ErrNetworkNotFound) {
			c.SendStatus(fiber.StatusNotFound)
			return nil
		}
		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return nil
		}
		nodeURL := fmt.Sprintf("http://localhost:%d/", nodeInstance.Config.Port)

		go saveOnRedis(ctx, c, nodeID)

		return proxy.Do(c, nodeURL)
	}
}

func saveOnRedis(ctx *Context, c *fiber.Ctx, nodeID string) {
	err := ctx.kvStore.PutMethodMetric(context.Background(), &storage.PutMethodMetricInput{
		NodeID: nodeID,
		Method: c.Method(),
	})
	if err != nil {
		log.Printf("error inserting metric for [%s-%s] %s\n", nodeID, c.Method(), err.Error())
		return
	}
}
