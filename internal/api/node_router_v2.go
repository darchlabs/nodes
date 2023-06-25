package api

import (
	"net/http"

	"github.com/darchlabs/backoffice/pkg/client"
	"github.com/darchlabs/backoffice/pkg/middleware"
	"github.com/darchlabs/nodes/internal/storage/instance"
	"github.com/gofiber/fiber/v2"
)

func routeV2Endpoints(ctx *Context) {
	// # Handlers
	postNewNodeV2Handler := &PostNewNodeV2Handler{
		instanceInsertQuery: instance.InsertQuery,
	}
	getNodesV2Handler := &GetNodesV2Handler{
		instanceSelectAllByUserIDQuery: instance.SelectAllByUserID,
	}
	deleteNodesV2Handler := &DeleteNodesV2Handler{
		instanceUpdateQuery:         instance.UpdateQuery,
		instanceSelectByUserIDQuery: instance.SelectByUserIDQuery,
	}

	// middlewares

	cl := client.New(&client.Config{
		Client:  http.DefaultClient,
		BaseURL: ctx.app.Config.BackofficeApiURL,
	})

	auth := middleware.NewAuth(cl)

	// # Route endpounts
	ctx.server.server.Get("/api/v2/health", handleFunc(
		ctx,
		func(_ *Context, _ *fiber.Ctx) (interface{}, int, error) {
			return map[string]string{"status": "running"}, fiber.StatusOK, nil
		},
	))

	ctx.server.server.Post("/api/v2/nodes", auth.Middleware, handleFunc(ctx, postNewNodeV2Handler.Invoke))
	ctx.server.server.Get("/api/v2/nodes", auth.Middleware, handleFunc(ctx, getNodesV2Handler.Invoke))
	ctx.server.server.Delete("/api/v2/nodes", auth.Middleware, handleFunc(ctx, deleteNodesV2Handler.Invoke))
}
