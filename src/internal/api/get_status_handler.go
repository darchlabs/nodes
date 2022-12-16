package api

import (
	"github.com/gofiber/fiber/v2"
)

type nodeStatus struct {
	ID              string `json:"id"`
	Chain           string `json:"chain"`
	Port            int    `json:"port"`
	FromBlockNumber int64  `json:"from_block_number"`
	Status          string `json:"status"`
}

type getStatusHandlerResponse struct {
	Nodes []*nodeStatus `json:"nodes"`
}

func getStatusHandler(ctx *Context, _ *fiber.Ctx) (interface{}, int, error) {
	nodeStatuses := make([]*nodeStatus, 0)

	nodeCommands := ctx.server.nodesManager.GetAll()

	for _, nodeCmd := range nodeCommands {
		nodeStatuses = append(nodeStatuses, &nodeStatus{
			ID:              nodeCmd.ID,
			Chain:           nodeCmd.Config.Network,
			Port:            nodeCmd.Config.Port,
			Status:          nodeCmd.Node.Status().String(),
			FromBlockNumber: nodeCmd.Config.FromBlockNumber,
		})
	}

	return &getStatusHandlerResponse{
		Nodes: nodeStatuses,
	}, fiber.StatusOK, nil
}
