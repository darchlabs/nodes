package api

import (
	"github.com/gofiber/fiber/v2"
)

type nodeStatus struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Chain           string `json:"chain"`
	Port            int    `json:"port"`
	FromBlockNumber int64  `json:"fromBlockNumber"`
	Status          string `json:"status"`
}

type getStatusHandlerResponse struct {
	Nodes []*nodeStatus `json:"nodes"`
}

func getStatusHandler(ctx *Context, _ *fiber.Ctx) (interface{}, int, error) {
	nodeStatuses := make([]*nodeStatus, 0)

	nodeInstances := ctx.server.nodesManager.GetAll()

	for _, instance := range nodeInstances {
		nodeStatuses = append(nodeStatuses, &nodeStatus{
			ID:              instance.ID,
			Chain:           instance.Config.Network,
			Name:            instance.Name,
			Port:            instance.Config.Port,
			Status:          instance.Node.Status().String(),
			FromBlockNumber: instance.Config.FromBlockNumber,
		})
	}

	return &getStatusHandlerResponse{
		Nodes: nodeStatuses,
	}, fiber.StatusOK, nil
}
