package api

import (
	"fmt"

	"github.com/darchlabs/nodes/src/internal/manager"
	"github.com/darchlabs/nodes/src/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type ServerConfig struct {
	Port              string
	Chain             string
	BootstrapNodesURL map[string]string
	Manager           *manager.Manager
}

type Server struct {
	server *fiber.App
	port   string

	nodesManager *manager.Manager
}

type Context struct {
	server *Server
	store  storage.DataStore
}

func NewServer(config *ServerConfig) *Server {

	return &Server{
		server:       fiber.New(),
		port:         config.Port,
		nodesManager: config.Manager,
	}
}

func (s *Server) Start(store storage.DataStore) error {
	go func() {
		ctx := &Context{
			server: s,
			store:  store,
		}
		// route endpoints
		routeNodeEndpoints("/api/v1/nodes", ctx)

		// proxy requests for node
		s.server.Post("/jsonrpc/:node_id", handleFunc(ctx, proxyRpcHandler))

		// sever listen
		fmt.Println("running")
		err := s.server.Listen(fmt.Sprintf(":%s", s.port))
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

type handler func(*Context, *fiber.Ctx) (interface{}, int, error)

func handleFunc(ctx *Context, fn handler) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		payload, statusCode, err := fn(ctx, c)
		if err != nil {
			return c.Status(statusCode).JSON(map[string]string{
				"error": err.Error(),
			})
		}

		if statusCode == statusAlreadyProxied {
			return nil
		}

		return c.Status(statusCode).JSON(payload)
	}
}
