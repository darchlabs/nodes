package api

import (
	"fmt"

	"github.com/darchlabs/nodes/src/internal/command"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type postNewNodeHandlerResponse struct {
	ID     string `json:"id"`
	Chain  string `json:"chain"`
	Port   int    `json:"port"`
	Status string `json:"status"`
}

func postNewNodeHandler(s *Server, _ *fiber.Ctx) (interface{}, int, error) {
	id := s.idGenerator()
	s.nodeConfig.Port++

	cmd := command.New(
		"ganache",
		"-p", fmt.Sprintf("%d", s.nodeConfig.Port),
		"--host", s.nodeConfig.Host,
		"--db", fmt.Sprintf("%s/%d", s.nodeConfig.BaseChainDataPath, len(s.nodesCommands)),
		"--fork", s.nodeConfig.BootsrapNodeURL,
	)

	err := cmd.StreamOutput(id)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: portNewNodeHandler cmd.StreamOutput error")
	}

	err = cmd.Start()
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: portNewNodeHandler cmd.Start error")
	}

	s.nodesCommands[id] = &nodeCommand{
		node: cmd,
		config: &NodeConfig{
			Host:              s.nodeConfig.Host,
			Port:              s.nodeConfig.Port,
			BaseChainDataPath: s.nodeConfig.DatabasePath,
			BootsrapNodeURL:   s.nodeConfig.BootsrapNodeURL,
		},
	}

	return &postNewNodeHandlerResponse{
		ID:     id,
		Chain:  s.chain,
		Port:   s.nodeConfig.Port,
		Status: cmd.Status().String(),
	}, fiber.StatusCreated, nil
}
