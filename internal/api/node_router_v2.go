package api

import (
	"github.com/darchlabs/nodes/internal/storage/instance"
	"github.com/gofiber/fiber/v2"
)

func routeV2Endpoints(ctx *Context) {
	// # Handlers
	postNewNodeV2Handler := &PostNewNodeV2Handler{
		instanceInsertQuery: instance.InsertQuery,
	}
	getNodesV2Handler := &GetNodesV2Handler{
		instanceSelectAllQuery: instance.SelectAllQuery,
	}
	deleteNodesV2Handler := &DeleteNodesV2Handler{
		instanceUpdateQuery: instance.UpdateQuery,
		instanceSelectQuery: instance.SelectQuery,
	}

	// # Route endpounts
	ctx.server.server.Get("/api/v2/health", handleFunc(
		ctx,
		func(_ *Context, _ *fiber.Ctx) (interface{}, int, error) {
			return map[string]string{"status": "running"}, fiber.StatusOK, nil
		},
	))

	ctx.server.server.Post("/api/v2/nodes", handleFunc(ctx, postNewNodeV2Handler.Invoke))
	ctx.server.server.Get("/api/v2/nodes", handleFunc(ctx, getNodesV2Handler.Invoke))
	ctx.server.server.Delete("/api/v2/nodes", handleFunc(ctx, deleteNodesV2Handler.Invoke))
}
