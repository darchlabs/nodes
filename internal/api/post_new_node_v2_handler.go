package api

import (
	"net/http"
	"time"

	"github.com/darchlabs/nodes/internal/manager"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type postNewNodev2HandlerRequest struct {
	Network string            `json:"network"`
	EnvVars map[string]string `json:"envVars"`
}

type PostNewNodev2HandlerResponse struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Network     string      `json:"network"`
	Environment string      `json:"environment,omitempty"`
	Port        int         `json:"port"`
	Status      string      `json:"status"`
	Artifacts   interface{} `json:"artifacts"`
	CreatedAt   time.Time   `json:"createdAt"`
}

func postNewNodeV2Handler(ctx *Context, c *fiber.Ctx) (interface{}, int, error) {
	var req postNewNodev2HandlerRequest
	err := c.BodyParser(&req)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: postNewNodeV2Handler c.BodyParser")
	}

	nodeInstance, err := ctx.server.nodesManager.DeployNewNode(&manager.CreatePodOptions{
		Network: req.Network,
		EnvVars: req.EnvVars,
	})
	if errors.Is(err, manager.ErrNodeNotFound) {
		return nil, http.StatusNotFound, nil
	}
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: postNewNodeV2Handler ctx.server.nodesManager.DeployNewNode")
	}

	return PostNewNodev2HandlerResponse{
		ID:          nodeInstance.ID,
		Name:        nodeInstance.Name,
		Network:     nodeInstance.Config.Network,
		Environment: nodeInstance.Config.Environment,
		Artifacts:   nodeInstance.Artifacts,
		Port:        nodeInstance.Config.Port,
		CreatedAt:   nodeInstance.Config.CreatedAt,
	}, fiber.StatusCreated, nil
}
