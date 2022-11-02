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
	Port          string
	CommandConfig *CommandConfig
}

type CommandConfig struct {
	Runner           string
	Host             string
	DatabasePath     string
	BootstrapNodeURL string
}

type Server struct {
	server *fiber.App
	port   string
	cmd    *command.Command
}

func NewServer(config *ServerConfig) *Server {
	return &Server{
		server: fiber.New(),
		port:   config.Port,
		cmd: command.New(
			config.CommandConfig.Runner,
			config.CommandConfig.Host,
			config.CommandConfig.DatabasePath,
			config.CommandConfig.BootstrapNodeURL,
		),
	}
}

func (s *Server) Start() error {
	err := s.cmd.Start()
	if err != nil {
		return err
	}
	log.Printf("Node is %s\n", s.cmd.Status())

	go func() {
		err := s.server.Listen(fmt.Sprintf(":%s", s.port))
		if err != nil {
			panic(err)
		}
	}()

	return nil
}
