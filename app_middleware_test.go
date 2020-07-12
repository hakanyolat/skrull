package skrull

import (
	"github.com/gofiber/fiber"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"testing"
)

func Test_syncCookiesFromResponse(t *testing.T){
	middleware := syncCookiesFromResponse([]string{"k"})

	app := fiber.New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	c := fasthttp.AcquireCookie()
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer app.ReleaseCtx(ctx)
	defer fasthttp.ReleaseCookie(c)
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	c.SetKey("k")
	c.SetValue("v")
	res.Header.SetCookie(c)

	mCtx := &MiddlewareContext{
		Ctx: ctx,
		Req: req,
		Res: res,
	}

	middleware(mCtx)
	assert.Equal(t, "k=v", string(mCtx.Ctx.Fasthttp.Response.Header.PeekCookie("k")))
}

func Test_syncCookiesToRequest(t *testing.T){
	middleware := syncCookiesToRequest([]string{"k"})

	app := fiber.New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer app.ReleaseCtx(ctx)
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	ctx.Fasthttp.Request.Header.SetCookie("k", "v")

	mCtx := &MiddlewareContext{
		Ctx: ctx,
		Req: req,
		Res: res,
	}

	middleware(mCtx)

	assert.Equal(t, "v", string(mCtx.Req.Header.Cookie("k")))
}
