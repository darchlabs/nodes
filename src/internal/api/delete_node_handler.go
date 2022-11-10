package api

import "github.com/gofiber/fiber/v2"

func deleteNodeHandler(s *Server, _ *fiber.Ctx) (interface{}, int, error) {
	return nil, fiber.StatusOK, nil
}
