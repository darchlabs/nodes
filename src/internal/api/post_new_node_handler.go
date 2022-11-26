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

func postNewNodeHandler(ctx *Context, _ *fiber.Ctx) (interface{}, int, error) {
	id := ctx.server.idGenerator()
	ctx.server.nodeConfig.Port++

	cmd := command.New(
		"ganache",
		"-p", fmt.Sprintf("%d", ctx.server.nodeConfig.Port),
		"--host", ctx.server.nodeConfig.Host,
		"--db", fmt.Sprintf("%s/%d", ctx.server.nodeConfig.BaseChainDataPath, len(ctx.server.nodesCommands)),
		"--fork", ctx.server.nodeConfig.BootsrapNodeURL,
	)

	err := cmd.StreamOutput(id)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: portNewNodeHandler cmd.StreamOutput error")
	}

	err = cmd.Start()
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "api: portNewNodeHandler cmd.Start error")
	}

	ctx.server.nodesCommands[id] = &nodeCommand{
		node: cmd,
		config: &NodeConfig{
			Host:              ctx.server.nodeConfig.Host,
			Port:              ctx.server.nodeConfig.Port,
			BaseChainDataPath: ctx.server.nodeConfig.DatabasePath,
			BootsrapNodeURL:   ctx.server.nodeConfig.BootsrapNodeURL,
		},
	}

	return &postNewNodeHandlerResponse{
		ID:     id,
		Chain:  ctx.server.chain,
		Port:   ctx.server.nodeConfig.Port,
		Status: cmd.Status().String(),
	}, fiber.StatusCreated, nil
}
