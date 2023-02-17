package api

import (
	"github.com/darchlabs/nodes/internal/command"
	"github.com/darchlabs/nodes/internal/manager"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

const (
	actionStart   = "start"
	actionStop    = "stop"
	actionRestart = "restart"
)

var nodeActions = map[string]bool{
	actionStart:   true,
	actionStop:    true,
	actionRestart: true,
}

type postActionHandlerRequest struct {
	Action string `json:"action"`
	NodeID string `json:"nodeId"`
}

type postActionHandlerResponse struct {
	ID     string `json:"id"`
	Chain  string `json:"chain"`
	Port   int    `json:"port"`
	Status string `json:"status"`
}

func postActionHandler(ctx *Context, c *fiber.Ctx) (interface{}, int, error) {
	var req postActionHandlerRequest
	err := c.BodyParser(&req)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: postActionHandler c.BodyParser error")
	}

	_, ok := nodeActions[req.Action]
	if !ok {
		return nil, fiber.StatusNotFound, errors.Wrap(ErrNotFound, "api: postActionHandler unrecognized status")
	}

	cmd, err := ctx.server.nodesManager.Get(req.NodeID)
	if errors.Is(err, manager.ErrNodeNotFound) {
		return nil, fiber.StatusNotFound, errors.Wrap(ErrNotFound, "api: postActionHandler ctx.server.nodesCommands unrecognized id")
	}
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: postActionHandler ctx.server.nodesManager.Get error")
	}

	nodeStatus := cmd.Node.Status()

	switch req.Action {
	case actionStart:
		if nodeStatus == command.StatusRunning {
			return nil, fiber.StatusOK, nil
		}

		err = cmd.Node.Start()
		if err != nil {
			return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: postActionHandler cmd.node.Start at starting error")
		}

	case actionStop:
		if nodeStatus == command.StatusStopped {
			return nil, fiber.StatusOK, nil
		}

		err = cmd.Node.Stop()
		if err != nil {
			return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: postActionHandler cmd.node.Stop at stopping error")
		}

	case actionRestart:
		err = cmd.Node.Stop()
		if err != nil {
			return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: postActionHandler cmd.node.Stop at restarting error")
		}

		node := cmd.Node.Clone()
		cmd.Node = node
		ctx.server.nodesManager.Put(cmd)
		err = cmd.Node.Start()
		if err != nil {
			return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: postActionHandler cmd.node.Start at restarting error")
		}
	}

	return nil, fiber.StatusOK, nil
}
