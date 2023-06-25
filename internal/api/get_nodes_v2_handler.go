package api

import (
	"time"

	"github.com/darchlabs/nodes/internal/storage/instance"
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
	instanceSelectAllByUserIDQuery instanceSelectAllByUserIDQuery
}

func (h *GetNodesV2Handler) Invoke(ctx *Context, c *fiber.Ctx) (interface{}, int, error) {
	userID, err := getUserIDFromRequestCtx(c)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: GetNodesV2Handler.Invoke getUserIDFromRequestCtx error")
	}

	payload, status, err := h.invoke(ctx, userID)
	if err != nil {
		return nil, status, errors.Wrap(err, "api: GetNodesV2Handler.Invoke h.invoke error")
	}

	return payload, status, nil
}

func (h *GetNodesV2Handler) invoke(ctx *Context, userID string) (interface{}, int, error) {
	records, err := h.instanceSelectAllByUserIDQuery(ctx.sqlStore, &instance.SelectAllByUserIDQuery{
		UserID: userID,
	})
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
