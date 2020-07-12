package skrull

import (
	"github.com/gofiber/fiber"
	"github.com/valyala/fasthttp"
)

type Middleware func(ctx *MiddlewareContext)

type MiddlewareContext struct {
	Ctx   *fiber.Ctx
	Req   *fasthttp.Request
	Res   *fasthttp.Response
	Store *MiddlewareStore
}
