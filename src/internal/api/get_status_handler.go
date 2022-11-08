package api

import (
	"github.com/gofiber/fiber/v2"
)

type nodeStatus struct {
	ID     string `json:"id"`
	Chain  string `json:"chain"`
	Status string `json:"status"`
}

type getStatusResponse struct {
	Nodes []*nodeStatus `json:"nodes"`
}

func getStatusHandler(s *Server, _ *fiber.Ctx) (interface{}, int, error) {
	nodeStatuses := make([]*nodeStatus, 0)

	for id, node := range s.nodes {
		nodeStatuses = append(nodeStatuses, &nodeStatus{
			ID:     id,
			Chain:  s.chain,
			Status: node.Status().String(),
		})
	}

	return &getStatusResponse{
		Nodes: nodeStatuses,
	}, fiber.StatusOK, nil
}
