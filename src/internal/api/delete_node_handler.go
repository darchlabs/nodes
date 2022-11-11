package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type deleteNodeHandlerRequest struct {
	NodeID string `json:"node_id"`
}

func deleteNodeHandler(s *Server, c *fiber.Ctx) (interface{}, int, error) {
	var req deleteNodeHandlerRequest
	err := c.BodyParser(&req)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: deleteNodeHandler c.BodyParser error")
	}

	cmd, ok := s.nodesCommands[req.NodeID]
	if !ok {
		return nil, fiber.StatusNotFound, errors.Wrap(err, "api: deleteNodeHandler s.nodeConfig unrecognized id")
	}

	err = cmd.node.Stop()
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: deleteNodeHandler node.cmd.Stop error")
	}

	delete(s.nodesCommands, req.NodeID)
	return nil, fiber.StatusOK, nil
}
