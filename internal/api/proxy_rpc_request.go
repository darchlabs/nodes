package api

import (
	"context"
	"fmt"
	"log"

	"github.com/darchlabs/nodes/internal/manager"
	"github.com/darchlabs/nodes/internal/storage"
	"github.com/darchlabs/nodes/internal/storage/instance"
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

type ProxyHandler struct {
	instanceSelectQuery instanceSelectQuery
}

func (h *ProxyHandler) invoke(ctx *Context) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		nodeID := c.Params("node_id")

		url, err := h.handleV1Search(nodeID, ctx)
		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}
		if url != "" {
			fmt.Println("~~~~~~> FORWARDED TO ", url)
			go saveOnRedis(ctx, c, nodeID)
			return proxy.Do(c, url)
		}

		url, err = h.handleV2Search(nodeID, ctx)
		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}
		if url == "" {
			c.SendStatus(fiber.StatusNotFound)
			return nil
		}

		return proxy.Do(c, url)
	}
}

func (h *ProxyHandler) handleV1Search(nodeID string, ctx *Context) (string, error) {
	nodeInstance, err := ctx.server.nodesManager.Get(nodeID)
	if errors.Is(err, manager.ErrNetworkNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://localhost:%d/", nodeInstance.Config.Port), nil
}

func (h *ProxyHandler) handleV2Search(nodeID string, ctx *Context) (string, error) {
	nodeInstance, err := h.instanceSelectQuery(ctx.sqlStore, &instance.SelectQueryInput{
		ID: nodeID,
	})
	if errors.Is(err, instance.ErrNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", nodeInstance.ServiceURL), nil
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
