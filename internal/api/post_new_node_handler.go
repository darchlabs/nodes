package api

import (
	"net/http"
	"time"

	"github.com/darchlabs/nodes/internal/manager"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type postNewNodeHandlerRequest struct {
	Network         string `json:"network"`
	FromBlockNumber int64  `json:"fromBlockNumber"`
}

type postNewNodeHandlerResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Chain     string    `json:"chain"`
	Port      int       `json:"port"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func postNewNodeHandler(ctx *Context, c *fiber.Ctx) (interface{}, int, error) {
	var req postNewNodeHandlerRequest
	err := c.BodyParser(&req)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: postNewNodeHandler c.BodyParser")
	}

	nodeInstance, err := ctx.server.nodesManager.CreateNode(&manager.CreateNodeConfig{
		Network:         req.Network,
		FromBlockNumber: req.FromBlockNumber,
	})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: postNewNodeHandler ctx.server.nodesManager.CreateNode error")
	}

	return &postNewNodeHandlerResponse{
		ID:        nodeInstance.ID,
		Name:      nodeInstance.Name,
		Chain:     req.Network,
		Port:      nodeInstance.Config.Port,
		Status:    nodeInstance.Node.Status().String(),
		CreatedAt: nodeInstance.Config.CreatedAt,
	}, fiber.StatusCreated, nil
}