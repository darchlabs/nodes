package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func routeNodeEndpoints(prefix string, ctx *Context) {
	ctx.server.server.Get("/api/health", handleFunc(ctx, func(_ *Context, _ *fiber.Ctx) (interface{}, int, error) {
		return map[string]string{"status": "runnin"}, fiber.StatusOK, nil
	}))

	ctx.server.server.Post(fmt.Sprintf("%s", prefix), handleFunc(ctx, postNewNodeHandler))
	ctx.server.server.Post(fmt.Sprintf("%s/actions", prefix), handleFunc(ctx, postActionHandler))
	ctx.server.server.Delete(fmt.Sprintf("%s", prefix), handleFunc(ctx, deleteNodeHandler))
	ctx.server.server.Get(fmt.Sprintf("%s/status", prefix), handleFunc(ctx, getStatusHandler))
	ctx.server.server.Get(fmt.Sprintf("%s/metrics", prefix), handleFunc(ctx, getNodesMetricsHandler))
}
