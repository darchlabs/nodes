package api

import (
	"fmt"
	"log"

	"github.com/darchlabs/nodes/src/internal/command"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type IDGenerator func() string

type ServerConfig struct {
	Port        string
	Chain       string
	IDGenerator IDGenerator
	Command     *command.Command
}

type Server struct {
	server       *fiber.App
	port         string
	cmd          *command.Command
	nodes        map[string]*command.Command
	idGenerator  IDGenerator
	masterNodeID string
	chain        string
}

func NewServer(config *ServerConfig) *Server {
	id := config.IDGenerator()

	return &Server{
		server:       fiber.New(),
		port:         config.Port,
		chain:        config.Chain,
		idGenerator:  config.IDGenerator,
		masterNodeID: id,
		nodes: map[string]*command.Command{
			id: config.Command,
		},
		cmd: config.Command,
	}
}

func (s *Server) Start() error {
	cmd, ok := s.nodes[s.masterNodeID]
	if !ok {
		return errors.New("api: Server.Start s.cmd.nodes not found")
	}

	err := cmd.Start(s.masterNodeID)
	if err != nil {
		return errors.Wrap(err, "api: Server.Start s.cmd.nodes not found")
	}

	log.Printf("Master %s-node is %s with id %s\n", s.chain, s.cmd.Status(), s.masterNodeID)

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
