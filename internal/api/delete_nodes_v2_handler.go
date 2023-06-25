package api

import (
	"net/http"
	"time"

	"github.com/darchlabs/nodes/internal/manager"
	"github.com/darchlabs/nodes/internal/storage/instance"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type DeleteNodesV2Handler struct {
	instanceUpdateQuery         instanceUpdateQuery
	instanceSelectByUserIDQuery instanceSelectByUserIDQuery
}

type deleteNodesV2HandlerRequest struct {
	ID     string `json:"id"`
	userID string `json:"-"`
}

func (h *DeleteNodesV2Handler) Invoke(ctx *Context, c *fiber.Ctx) (interface{}, int, error) {
	var req deleteNodesV2HandlerRequest
	err := c.BodyParser(&req)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: DeleteNodesV2HandlerRequest.Invoke c.BodyParser")
	}

	req.userID, err = getUserIDFromRequestCtx(c)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: DeleteNodesV2Handler.Invoke getUserIDFromRequestCtx error")
	}

	_, status, err := h.invoke(ctx, &req)
	if err != nil {
		return nil, status, errors.Wrap(err, "api: DeleteNodesV2Handler.Invoke h.invoke error")
	}
	return nil, status, nil
}

func (h *DeleteNodesV2Handler) invoke(ctx *Context, req *deleteNodesV2HandlerRequest) (interface{}, int, error) {
	// select instance
	instanceRecord, err := h.instanceSelectByUserIDQuery(ctx.sqlStore, &instance.SelectByUserIDQueryInput{
		ID:     req.ID,
		UserID: req.userID,
	})
	if errors.Is(err, instance.ErrNotFound) {
		return nil, fiber.StatusNotFound, nil
	}
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: DeleteNodesV2HandlerRequest.invoke h.instanceSelectByUSerIDQuery error")
	}

	// delete artifacts
	err = ctx.nodeManager.DeleteNode(&manager.NodeInstance{
		ID:   instanceRecord.ID,
		Name: instanceRecord.Name,
		Config: &manager.NodeConfig{
			Network:     instanceRecord.Network,
			Environment: instanceRecord.Environment,
		},
		Artifacts: &manager.Artifacts{
			Deployments: instanceRecord.Artifacts.Deployments,
			Pods:        instanceRecord.Artifacts.Pods,
			Services:    instanceRecord.Artifacts.Services,
		},
	})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: DeleteNodesV2HandlerRequest.invoke ctx.nodeManager.DeleteArtifacts error")
	}

	// update instance
	now := time.Now()
	err = h.instanceUpdateQuery(ctx.sqlStore, &instance.UpdateQueryInput{
		ID:        req.ID,
		UserID:    req.userID,
		DeletedAt: &now,
	})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: DeleteNodesV2HandlerRequest.invoke h.instanceUpdateQuery error")
	}

	return nil, fiber.StatusNoContent, nil
}
