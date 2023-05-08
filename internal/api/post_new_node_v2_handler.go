package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/darchlabs/nodes/internal/manager"
	"github.com/darchlabs/nodes/internal/storage/instance"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type PostNewNodeV2Handler struct {
	instanceInsertQuery instanceInsertQuery
}

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

func (h *PostNewNodeV2Handler) Invoke(ctx *Context, c *fiber.Ctx) (interface{}, int, error) {
	var req *postNewNodev2HandlerRequest
	err := c.BodyParser(&req)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: postNewNodeV2Handler c.BodyParser")
	}

	payload, status, err := h.invoke(ctx, req)
	if err != nil {
		return nil, status, errors.Wrap(err, "api: PostNewNodeV2Handler.Invoke h.invoke error")
	}

	return payload, status, nil
}

func (h *PostNewNodeV2Handler) invoke(ctx *Context, req *postNewNodev2HandlerRequest) (interface{}, int, error) {
	nodeInstance, err := ctx.nodeManager.DeployNewNode(&manager.CreateDeploymentOptions{
		Network: req.Network,
		EnvVars: req.EnvVars,
	})
	if errors.Is(err, manager.ErrNetworkNotFound) {
		return nil, http.StatusNotFound, nil
	}
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: PostNewNodeV2Handler.invoke h.nodesManager.DeployNewNode")
	}

	instanceRecord := &instance.Record{
		ID:          nodeInstance.ID,
		Network:     nodeInstance.Config.Network,
		Environment: nodeInstance.Config.Environment,
<<<<<<< HEAD
		ServiceURL:  fmt.Sprintf("http://%s:%d", nodeInstance.Name, nodeInstance.Config.Port),
		Name:        nodeInstance.Name,
		CreatedAt:   time.Now(),
		Artifacts: &instance.Artifacts{
			Pods:        nodeInstance.Artifacts.Pods,
			Deployments: nodeInstance.Artifacts.Deployments,
			Services:    nodeInstance.Artifacts.Services,
		},
=======
		Name:        nodeInstance.Name,
		CreatedAt:   time.Now(),
>>>>>>> c30c8a7 (Update api to support sql database)
	}

	err = h.instanceInsertQuery(ctx.sqlStore, instanceRecord)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: PostNewNodeV2Handler.invoke instance.InsertQuery error")
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
