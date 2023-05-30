package api

import (
	"github.com/darchlabs/nodes/internal/manager"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type deleteNodeHandlerV2Request struct {
	NodeID string `json:"nodeId"`
}

func deleteNodeHandler(ctx *Context, c *fiber.Ctx) (interface{}, int, error) {
	var req deleteNodeHandlerV2Request
	err := c.BodyParser(&req)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: deleteNodeHandler c.BodyParser error")
	}

	cmd, err := ctx.server.nodesManager.Get(req.NodeID)
	if errors.Is(err, manager.ErrNetworkNotFound) {
		return nil, fiber.StatusNotFound, errors.Wrap(ErrNotFound, "api: deleteNodeHandler ctx.server.nodeConfig unrecognized id")
	}
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: deleteNodeHandler ctx.server.nodeConfig unrecognized id")
	}

	err = cmd.Node.Stop()
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: deleteNodeHandler node.cmd.Stop error")
	}

	ctx.server.nodesManager.Delete(cmd.ID)
	return nil, fiber.StatusOK, nil
}
