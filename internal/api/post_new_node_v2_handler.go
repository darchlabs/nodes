package api

import (
	"net/http"

	"github.com/darchlabs/nodes/internal/manager"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type postNewNodev2HandlerRequest struct {
	Image   string            `json:"image"`
	EnvVars map[string]string `json:"env_vars,envVars"`
}

type PostNewNodev2HandlerResponse struct {
	Name     map[string]string `json:"name"`
	NodeType string            `json:"node_type,nodeType"`
}

func postNewNodeV2Handler(ctx *Context, c *fiber.Ctx) (interface{}, int, error) {
	var req postNewNodev2HandlerRequest
	err := c.BodyParser(&req)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: postnewnodev2handler c.BodyParser")
	}

	err = ctx.server.nodesManager.CreatePod(&manager.CreatePodOptions{
		Image:   req.Image,
		EnvVars: req.EnvVars,
	})

	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: postnewnodev2handler ctx.server.nodesManager.CreatePod")
	}

	return nil, fiber.StatusCreated, nil
}
