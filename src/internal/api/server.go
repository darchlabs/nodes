package api

import (
	"fmt"
	"log"

	"github.com/darchlabs/nodes/src/internal/command"
	"github.com/gofiber/fiber/v2"
)

type Context struct {
	cmd *command.Command
}

type ServerConfig struct {
	Port    string
	Command *command.Command
}

type Server struct {
	server *fiber.App // TODO: Replace with real server instance
	port   string
	cmd    *command.Command
}

func NewServer(config *ServerConfig) *Server {
	server := fiber.New()

	return &Server{
		server: server,
		port:   config.Port,
	}
}

func (s *Server) Start() error {
	err := s.cmd.StreamOutput()
	err = s.cmd.Start()
	if err != nil {
		return err
	}

	log.Println(s.cmd.Status())

	go func() { // TODO: replace this with server running
		s.server.Listen(fmt.Sprintf(":%s", s.port))
	}()

	return nil
}
