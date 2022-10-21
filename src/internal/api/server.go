package api

import (
	"log"

	"github.com/darchlabs/nodes/src/internal/command"
	"github.com/gin-gonic/gin"
)

type Context struct {
	cmd *command.Command
}

type ServerConfig struct {
	Port    string
	Command *command.Command
}

type Server struct {
	server *gin.Engine // TODO: Replace with real server instance
	port   string
	cmd    *command.Command
}

func NewServer(config *ServerConfig) *Server {
	return &Server{
		//server: gin.Default(),
		port: config.Port,
		cmd:  config.Command,
	}
}

func (s *Server) Start() error {
	err := s.cmd.Start()
	if err != nil {
		return err
	}

	log.Println(s.cmd.Slug())
	log.Println(s.cmd.Status())

	go func() { // TODO: replace this with server running
		//fmt.Println(s.server)
		//s.server.Run(fmt.Sprintf(":%s", s.port))
	}()

	return nil
}
