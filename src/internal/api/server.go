package api

import (
	"fmt"

	"github.com/darchlabs/nodes/src/internal/command"
	"github.com/gin-gonic/gin"
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
	return &Server{
		server: gin.Default(),
		port:   config.Port,
		cmd:    config.Command,
	}
}

func (s *Server) Start() error {
	//log.Println(s.cmd.Status())

	go func() { // TODO: replace this with server running
		s.server.Listen(fmt.Sprintf(":%s", s.port))
	}()

	return nil
}
