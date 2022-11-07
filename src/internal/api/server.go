package api

import (
	"fmt"
	"log"

	"github.com/darchlabs/nodes/src/internal/command"
	"github.com/gofiber/fiber/v2"
)

type ServerConfig struct {
	Port    string
	Chain   string
	Command *command.Command
}

type Server struct {
	server *fiber.App
	cmd    *command.Command
	port   string
	chain  string
}

func NewServer(config *ServerConfig) *Server {
	return &Server{
		server: fiber.New(),
		port:   config.Port,
		chain:  config.Chain,
		cmd:    config.Command,
	}
}

func (s *Server) Start() error {
	err := s.cmd.Start()
	if err != nil {
		return err
	}
	log.Printf("Node is %s\n", s.cmd.Status())

	go func() {
		// route endpoints
		routeNodeEndpoints("/nodes", s)

		// sever listen
		err := s.server.Listen(fmt.Sprintf(":%s", s.port))
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

type handler func(*Server, *fiber.Ctx) (interface{}, int, error)

func handleFunc(s *Server, fn handler) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		payload, statusCode, err := fn(s, c)
		if err != nil {
			return c.Status(statusCode).JSON(map[string]string{
				"error": err.Error(),
			})
		}

		return c.Status(statusCode).JSON(payload)
	}
}
