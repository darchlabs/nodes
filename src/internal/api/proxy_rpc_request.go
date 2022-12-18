package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/darchlabs/nodes/src/internal/manager"
	"github.com/darchlabs/nodes/src/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

var statusAlreadyProxied int = 1000

type proxyRpcHandlerRequest struct {
	JSONRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      interface{}   `json:"id"`
}

func proxyRpcHandler(ctx *Context, c *fiber.Ctx) (interface{}, int, error) {
	// get path params
	nodeID := c.Params("node_id")
	log.Println("request forwarded for for node", nodeID)

	cmd, err := ctx.server.nodesManager.Get(nodeID)
	if errors.Is(err, manager.ErrNodeNotFound) {
		return nil, fiber.StatusNotFound, errors.Wrap(ErrNotFound, "api: proxyRpcHandler unrecognized node_id")
	}
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(ErrNotFound, "api: proxyRpcHandler ctx.server.nodesManager.Get error")
	}

	var req proxyRpcHandlerRequest
	err = c.BodyParser(&req)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: proxyRpcHandler c.BodyParser error")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: proxyRpcHandler json.Marshal error")
	}

	go func() {
		err := ctx.store.PutMethodMetric(context.Background(), &storage.PutMethodMetricInput{
			NodeID: nodeID,
			Method: req.Method,
		})
		if err != nil {
			log.Printf("error inserting metric for [%s-%s]\n", nodeID, req.Method)
		}
	}()

	nodeURL := fmt.Sprintf("http://0.0.0.0:%d/", cmd.Config.Port)
	request, err := http.NewRequest(http.MethodPost, nodeURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: proxyRpcHandler http.NewRequest error")
	}

	headers := c.GetReqHeaders()

	for k, v := range headers {
		request.Header.Add(k, v)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: proxyRpcHandler http.DefaultClient.Do error")
	}

	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: proxyRpcHandler ioutil.ReadAll error")
	}

	return nil, statusAlreadyProxied, c.Send(bodyBytes)
}
