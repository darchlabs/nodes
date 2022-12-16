package api

import (
	"net/http"

	"github.com/darchlabs/nodes/src/internal/manager"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type postNewNodeHandlerRequest struct {
	Network         string `json:"network"`
	FromBlockNumber int64  `json:"from_block_number"`
}

type postNewNodeHandlerResponse struct {
	ID     string `json:"id"`
	Chain  string `json:"chain"`
	Port   int    `json:"port"`
	Status string `json:"status"`
}

func postNewNodeHandler(ctx *Context, c *fiber.Ctx) (interface{}, int, error) {
	var req postNewNodeHandlerRequest
	err := c.BodyParser(&req)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: postNewNodeHandler c.BodyParser")
	}

	nodeCmd, err := ctx.server.nodesManager.CreateNode(&manager.CreateNodeConfig{
		Network:         req.Network,
		FromBlockNumber: req.FromBlockNumber,
	})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: postNewNodeHandler ctx.server.nodesManager.CreateNode error")
	}

	return &postNewNodeHandlerResponse{
		ID:     nodeCmd.ID,
		Chain:  req.Network,
		Port:   nodeCmd.Config.Port,
		Status: nodeCmd.Node.Status().String(),
	}, fiber.StatusCreated, nil
}
