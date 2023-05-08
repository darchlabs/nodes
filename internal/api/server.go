package api

import (
	"fmt"
	"os"

	"github.com/darchlabs/nodes/internal/application"
	"github.com/darchlabs/nodes/internal/manager"
	"github.com/darchlabs/nodes/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type ServerConfig struct {
	Port              string
	Chain             string
	BootstrapNodesURL map[string]string
	Manager           *manager.Manager
	App               *application.App
}

type Server struct {
	server *fiber.App
	app    *application.App
	port   string

	nodesManager *manager.Manager
}

type Context struct {
	// structs
	server *Server
	app    *application.App

	// interfaces
	nodeManager nodeManager
	kvStore     storage.KeyValue
	sqlStore    storage.SQL
}

func NewServer(config *ServerConfig) *Server {

	server := fiber.New()
	server.Use(logger.New())
	server.Use(logger.New(logger.Config{
		Format:     "[${ip}]:${port} ${status} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		Output:     os.Stdout,
	}))

	return &Server{
		server:       server,
		port:         config.Port,
		nodesManager: config.Manager,
	}
}

func (s *Server) Start(app *application.App) error {
	go func() {
		ctx := &Context{
			server:   s,
			kvStore:  app.KeyValueStore,
			sqlStore: app.SqlStore,
		}
		// route endpoints
		routeNodeEndpoints("/api/v1/nodes", ctx)
		routeV2Endpoints(ctx)

		// proxy requests for node
		s.server.All("jsonrpc/:node_id", proxyFunc(ctx))

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
