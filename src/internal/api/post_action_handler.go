package api

import (
	"github.com/darchlabs/nodes/src/internal/command"
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
	NodeID string `json:"node_id"`
}

type postActionHandlerResponse struct {
	ID     string `json:"id"`
	Chain  string `json:"chain"`
	Port   int    `json:"port"`
	Status string `json:"status"`
}

func postActionHandler(s *Server, c *fiber.Ctx) (interface{}, int, error) {
	var req postActionHandlerRequest
	err := c.BodyParser(&req)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: postActionHandler c.BodyParser error")
	}

	_, ok := nodeActions[req.Action]
	if !ok {
		return nil, fiber.StatusNotFound, errors.Wrap(ErrNotFound, "api: postActionHandler unrecognized status")
	}

	cmd, ok := s.nodesCommands[req.NodeID]
	if !ok {
		return nil, fiber.StatusNotFound, errors.Wrap(ErrNotFound, "api: postActionHandler s.nodesCommands unrecognized id")
	}

	nodeStatus := cmd.node.Status()

	switch req.Action {
	case actionStart:
		if nodeStatus == command.StatusRunning {
			return nil, fiber.StatusOK, nil
		}

		err = cmd.node.Start()
		if err != nil {
			return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: postActionHandler cmd.node.Start at starting error")
		}

	case actionStop:
		if nodeStatus == command.StatusStopped {
			return nil, fiber.StatusOK, nil
		}

		err = cmd.node.Stop()
		if err != nil {
			return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: postActionHandler cmd.node.Stop at stopping error")
		}

	case actionRestart:
		err = cmd.node.Stop()
		if err != nil {
			return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: postActionHandler cmd.node.Stop at restarting error")
		}

		s.nodesCommands[req.NodeID].node = cmd.node.Clone()
		err = s.nodesCommands[req.NodeID].node.Start()
		if err != nil {
			return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: postActionHandler cmd.node.Start at restarting error")
		}
	}

	return nil, fiber.StatusOK, nil
}
