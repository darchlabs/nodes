package api

import (
	"github.com/gofiber/fiber/v2"
)

type nodeStatus struct {
	ID     string `json:"id"`
	Chain  string `json:"chain"`
	Port   int    `json:"port"`
	Status string `json:"status"`
}

type getStatusHandlerResponse struct {
	Nodes []*nodeStatus `json:"nodes"`
}

func getStatusHandler(ctx *Context, _ *fiber.Ctx) (interface{}, int, error) {
	nodeStatuses := make([]*nodeStatus, 0)

	for id, cmd := range ctx.server.nodesCommands {
		nodeStatuses = append(nodeStatuses, &nodeStatus{
			ID:     id,
			Chain:  ctx.server.chain,
			Port:   cmd.config.Port,
			Status: cmd.node.Status().String(),
		})
	}

	return &getStatusHandlerResponse{
		Nodes: nodeStatuses,
	}, fiber.StatusOK, nil
}
