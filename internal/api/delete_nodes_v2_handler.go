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
	instanceUpdateQuery instanceUpdateQuery
	instanceSelectQuery instanceSelectQuery
}

type deleteNodesV2HandlerRequest struct {
	ID string
}

func (h *DeleteNodesV2Handler) Invoke(ctx *Context, c *fiber.Ctx) (interface{}, int, error) {
	var req deleteNodesV2HandlerRequest
	err := c.BodyParser(&req)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: DeleteNodesV2HandlerRequest.Invoke c.BodyParser")
	}

	_, status, err := h.invoke(ctx, &req)
	if err != nil {
		return nil, status, errors.Wrap(err, "api: DeleteNodesV2Handler.Invoke h.invoke error")
	}
	return nil, status, nil
}

func (h *DeleteNodesV2Handler) invoke(ctx *Context, req *deleteNodesV2HandlerRequest) (interface{}, int, error) {
	// select instance
	instanceRecord, err := h.instanceSelectQuery(ctx.sqlStore, &instance.SelectQueryInput{
		ID: req.ID,
	})
	if errors.Is(err, instance.ErrNotFound) {
		return nil, fiber.StatusNotFound, nil
	}
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: DeleteNodesV2HandlerRequest.invoke h.instanceSelectQuery error")
	}

	// delete artifacts
	err = ctx.nodeManager.DeleteArtifacts(&manager.Artifacts{
		Deployments: instanceRecord.Artifacts.Deployments,
		Pods:        instanceRecord.Artifacts.Pods,
		Services:    instanceRecord.Artifacts.Services,
	})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "api: DeleteNodesV2HandlerRequest.invoke ctx.nodeManager.DeleteArtifacts error")
	}

	// update instance
	now := time.Now()
	err = h.instanceUpdateQuery(ctx.sqlStore, &instance.UpdateQueryInput{
		ID:        req.ID,
		DeletedAt: &now,
	})

	return nil, fiber.StatusNoContent, nil
}
