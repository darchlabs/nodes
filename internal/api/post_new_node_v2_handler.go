package api

import (
	"net/http"
	"time"

	"github.com/darchlabs/nodes/internal/manager"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type postNewNodev2HandlerRequest struct {
	Image   string            `json:"image"`
	EnvVars map[string]string `json:"envVars"`
}

type PostNewNodev2HandlerResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Chain     string    `json:"chain"`
	Port      int       `json:"port"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func postNewNodeV2Handler(ctx *Context, c *fiber.Ctx) (interface{}, int, error) {
	var req postNewNodev2HandlerRequest
	err := c.BodyParser(&req)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: postnewnodev2handler c.BodyParser")
	}

	nodeInstance, err := ctx.server.nodesManager.CreatePod(&manager.CreatePodOptions{
		Image:   req.Image,
		EnvVars: req.EnvVars,
	})

	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: postnewnodev2handler ctx.server.nodesManager.CreatePod")
	}

	return PostNewNodev2HandlerResponse{
		ID:   nodeInstance.ID,
		Name: nodeInstance.Name,
	}, fiber.StatusCreated, nil
}
