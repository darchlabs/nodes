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
	NodeConfig  *NodeConfig
}

type NodeConfig struct {
	Host              string
	Port              int
	DatabasePath      string
	BaseChainDataPath string
	BootsrapNodeURL   string
}

type Server struct {
	server      *fiber.App
	idGenerator IDGenerator
	port        string

	nodeConfig    NodeConfig
	nodesCommands map[string]*nodeCommand
	masterNodeID  string
	chain         string
}

type nodeCommand struct {
	node   *command.Command
	config *NodeConfig
}

func NewServer(config *ServerConfig) *Server {
	id := config.IDGenerator()
	cmd := command.New(
		"ganache",
		"-p", fmt.Sprintf("%d", config.NodeConfig.Port),
		"--host", config.NodeConfig.Host,
		"--db", config.NodeConfig.BaseChainDataPath,
		"--fork", config.NodeConfig.BootsrapNodeURL,
	)

	return &Server{
		server:       fiber.New(),
		idGenerator:  config.IDGenerator,
		port:         config.Port,
		chain:        config.Chain,
		nodeConfig:   *config.NodeConfig,
		masterNodeID: id,
		nodesCommands: map[string]*nodeCommand{
			id: &nodeCommand{
				node:   cmd,
				config: config.NodeConfig,
			},
		},
	}
}

func (s *Server) Start() error {
	cmd, ok := s.nodesCommands[s.masterNodeID]
	if !ok {
		return errors.New("api: Server.Start s.cmd.nodes not found")
	}

	err := cmd.node.StreamOutput(s.masterNodeID)
	if err != nil {
		return errors.Wrap(err, "api: Server.Start cmd.node.StreamOutput error")
	}

	err = cmd.node.Start()
	if err != nil {
		return errors.Wrap(err, "api: Server.Start cmd.node.Start error")
	}

	log.Printf("Master %s-node is %s with id %s\n", s.chain, cmd.node.Status(), s.masterNodeID)

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
