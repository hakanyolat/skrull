package skrull

import (
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/gofiber/recover"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

var httpMethods = []string{
	fasthttp.MethodPost,
	fasthttp.MethodConnect,
	fasthttp.MethodDelete,
	fasthttp.MethodGet,
	fasthttp.MethodHead,
	fasthttp.MethodOptions,
	fasthttp.MethodPatch,
	fasthttp.MethodPut,
	fasthttp.MethodTrace,
}

type App struct {
	Url                    string
	address                interface{}
	logger                 *zap.Logger
	http                   *fiber.App
	ignores                []ignore
	fragments              []Fragment
	requestMiddlewareList  []Middleware
	responseMiddlewareList []Middleware
	miners                 []Miner
	cookieTable            []string
}

// NewApp returns new App instance
func NewApp(url string) *App {
	return &App{Url: url}
}

// Replace adds fragment(s) to replace queue
func (app *App) Replace(fragments ...Fragment) {
	app.fragments = append(app.fragments, fragments...)
}

// Craft
func (app *App) Craft(miner ...Miner) {
	app.miners = append(app.miners, miner...)
}

// RequestMiddleware adds middleware to before fake request
func (app *App) RequestMiddleware(middleware Middleware) {
	app.requestMiddlewareList = append(app.requestMiddlewareList, middleware)
}

// ResponseMiddleware adds middleware to after fake request
func (app *App) ResponseMiddleware(middleware Middleware) {
	app.responseMiddlewareList = append(app.responseMiddlewareList, middleware)
}

// Pass ignores content rendering for specific method and path
// The request is completely simulated via a fake site
func (app *App) Pass(method, path string) {
	app.ignores = append(app.ignores, ignore{
		Method: method,
		Path:   path,
	})
}

// PassAll passes path for all http methods
func (app *App) PassAll(path string) {
	for _, method := range httpMethods {
		app.Pass(method, path)
	}
}

// SyncCookie synchronizes all cookies between main request and fake request
func (app *App) SyncCookie(cookieName string) {
	app.cookieTable = append(app.cookieTable, cookieName)
}

// Listen creates and configures an http app
func (app *App) Listen(address interface{}) {
	app.address = address
	app.logger, _ = zap.NewProduction()
	app.http = fiber.New()
	app.http.Use(recover.New())
	app.http.All("/*", app.handle)
}

// Run runs the application
func (app *App) Run() {
	defer app.logger.Sync()

	if app.http == nil {
		panic("There is no http application listening. Please listen with Listen(address interface{}) first.")
	}

	app.RequestMiddleware(syncCookiesToRequest(app.cookieTable))
	app.ResponseMiddleware(syncCookiesFromResponse(app.cookieTable))
	go func() {
		app.logger.Fatal(app.http.Listen(app.address).Error())
	}()
	configureGracefulShutdown(app.logger, app.http)
}

// handle Decides what to do with the request
func (app *App) handle(ctx *fiber.Ctx) {
	var err error
	var body string

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	middlewareContext := &MiddlewareContext{
		Ctx:   ctx,
		Req:   req,
		Res:   res,
		Store: &MiddlewareStore{},
	}

	req.Header.SetMethod(ctx.Method())
	relativePath := createFullRelativePath(ctx.Path(), ctx.Fasthttp.QueryArgs().String())
	req.SetRequestURI(createFullPath(app.Url, relativePath))

	app.handleRequestMiddleware(middlewareContext)
	err = fasthttp.Do(req, res)

	if err != nil {
		app.logger.Error(err.Error())
	}

	statusCode := res.Header.StatusCode()
	if !statusIsMoved(statusCode) {
		body = string(res.Body())
		contentType := res.Header.ContentType()
		ctx.Set("Content-type", string(contentType))

		if !app.isIgnored(ctx.Method(), ctx.Path()) {
			app.handleResponseMiddleware(middlewareContext)
			body = app.mountBody(middlewareContext)
		}

		ctx.Write(body)
		return
	}

	app.handleResponseMiddleware(middlewareContext)
	if location := string(res.Header.Peek("Location")); location != "" {
		newLocation := strings.Replace(location, app.Url, "", -1)
		if newLocation == "" {
			newLocation = "/"
		}
		ctx.Redirect(newLocation)
	} else {
		s := fmt.Sprintf("Bad Request: Moved status(%d) with empty location", statusCode)
		ctx.Write(s)
		ctx.Set("Content-type", "text/html")
		ctx.Status(fasthttp.StatusBadRequest)
		app.logger.Error(s)
	}
}

// isIgnored decides that the request will be passed
func (app *App) isIgnored(method, path string) bool {
	for _, ignore := range app.ignores {
		if matched, _ := regexp.Match(ignore.Path, []byte(path)); matched && ignore.Method == method {
			return true
		}
	}
	return false
}

// handleRequestMiddleware handles middleware to before fake request
func (app *App) handleRequestMiddleware(ctx *MiddlewareContext) {
	for _, m := range app.requestMiddlewareList {
		if !app.isIgnored(ctx.Ctx.Method(), ctx.Ctx.Path()) {
			m(ctx)
		}
	}
}

// handleResponseMiddleware handles middleware to after fake request
func (app *App) handleResponseMiddleware(ctx *MiddlewareContext) {
	for _, m := range app.responseMiddlewareList {
		if !app.isIgnored(ctx.Ctx.Method(), ctx.Ctx.Path()) {
			m(ctx)
		}
	}
}

// mountBody produces content to be displayed as a result
func (app *App) mountBody(ctx *MiddlewareContext) string {
	parser := NewHtmlParser(ctx.Res.Body())

	ch := make(chan RenderedFragment)
	for _, fragment := range app.fragments {
		if fragment.Available(parser) {
			go fragment.Render(*ctx, ch, app.logger)
		}
	}

	for range app.fragments {
		res := <-ch
		parser.Replace(res)
	}

	for _, miner := range app.miners {
		miner(parser)
	}

	newHtml, err := parser.Document.Html()
	if err != nil {
		app.logger.Error(err.Error())
	}

	return newHtml
}
