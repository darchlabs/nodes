package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type nodeV2 struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Network     string      `json:"network"`
	Environment string      `json:"environment,omitempty"`
	Port        int         `json:"port"`
	Status      string      `json:"status"`
	Artifacts   interface{} `json:"artifacts"`
	CreatedAt   time.Time   `json:"createdAt"`
}

type getNodesV2HandlerResponse struct {
	Nodes []*nodeV2 `json:"nodes"`
}

type GetNodesV2Handler struct {
	instanceSelectAllQuery instanceSelectAllQuery
}

func (h *GetNodesV2Handler) Invoke(ctx *Context, c *fiber.Ctx) (interface{}, int, error) {
	payload, status, err := h.invoke(ctx)
	if err != nil {
		return nil, status, errors.Wrap(err, "api: GetNodesV2Handler.Invoke h.invoke error")
	}

	return payload, status, nil
}

func (h *GetNodesV2Handler) invoke(ctx *Context) (interface{}, int, error) {
	records, err := h.instanceSelectAllQuery(ctx.sqlStore)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "manager: GetNodesV2Handler.invoke h.instanceSelectAllQuery")
	}

	nodes := make([]*nodeV2, 0)
	for _, r := range records {
		nodes = append(nodes, &nodeV2{
			ID:          r.ID,
			Name:        r.Name,
			Network:     r.Network,
			Environment: r.Environment,
			Artifacts:   r.Artifacts,
			Status:      "running", // update status with real status
			CreatedAt:   r.CreatedAt,
		})
	}

	return &getNodesV2HandlerResponse{
		Nodes: nodes,
	}, fiber.StatusOK, nil
}
