package api

import (
	"fmt"
	"log"

	"github.com/darchlabs/nodes/src/internal/command"
	"github.com/gin-gonic/gin"
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
	server *gin.Engine
	port   string
	cmd    *command.Command
}

func NewServer(config *ServerConfig) *Server {
	return &Server{
		server: gin.Default(),
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
	go func() { // TODO: replace this with server running
		fmt.Println("Running server")
		err := s.server.Run(fmt.Sprintf(":%s", s.port))
		fmt.Println("ERROR ON SERVER RUNNING ", err)
	}()

	err := s.cmd.Start()
	if err != nil {
		return err
	}

	log.Println(s.cmd.Slug())
	log.Println(s.cmd.Status())

	return nil
}
