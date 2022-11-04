package api

import (
	"github.com/gofiber/fiber/v2"
)

type getStatusResponse struct {
	Chain  string `json:"chain"`
	Status string `json:"status"`
}

func getStatusHandler(s *Server, _ *fiber.Ctx) (interface{}, int, error) {
	status := s.cmd.Status()
	return &getStatusResponse{
		Chain:  s.chain,
		Status: status.String(),
	}, fiber.StatusOK, nil
}
