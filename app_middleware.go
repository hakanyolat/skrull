package skrull

import (
	"github.com/valyala/fasthttp"
)

// syncCookiesFromResponse ...
func syncCookiesFromResponse(cookieTable []string) Middleware {
	return func(ctx *MiddlewareContext) {
		ctx.Res.Header.VisitAllCookie(func(key, value []byte) {
			c := fasthttp.AcquireCookie()
			defer fasthttp.ReleaseCookie(c)
			cKey := string(key)
			for _, cookieName := range cookieTable {
				if cKey == cookieName {
					_ = c.ParseBytes(value)
					c.SetDomain("")
					ctx.Ctx.Fasthttp.Response.Header.SetCookie(c)
				}
			}
		})
	}
}

// syncCookiesToRequest ...
func syncCookiesToRequest(cookieTable []string) Middleware {
	return func(ctx *MiddlewareContext) {
		ctx.Ctx.Fasthttp.Request.Header.VisitAllCookie(func(key, value []byte) {
			c := fasthttp.AcquireCookie()
			defer fasthttp.ReleaseCookie(c)
			cKey := string(key)
			for _, cookieName := range cookieTable {
				if cKey == cookieName {
					_ = c.ParseBytes(value)
					ctx.Req.Header.SetCookie(cKey, string(c.Value()))
				}
			}
		})
	}
}