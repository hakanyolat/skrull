package skrull

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Fragment struct {
	Name     string
	Url      string
	Selector string
}

type RenderedFragment struct {
	Fragment
	Result string
}

// Select ...
func (f Fragment) Select(parser *HtmlParser) *goquery.Selection {
	return parser.Document.Find(f.Selector)
}

// Available ...
func (f Fragment) Available(parser *HtmlParser) bool {
	return parser.Document.Contains(f.Select(parser).Get(0))
}

// Render ...
func (f Fragment) Render(ctx MiddlewareContext, ch chan RenderedFragment, logger *zap.Logger) {
	path := ctx.Ctx.Path()
	args := ctx.Ctx.Fasthttp.QueryArgs().String()
	fullPath := createFullPath(f.Url, createFullRelativePath(path, args))

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)
	req.SetRequestURI(fullPath)
	err := fasthttp.Do(req, res)

	if err != nil {
		logger.Error(err.Error())
	}

	ch <- RenderedFragment{
		Fragment: f,
		Result:   string(res.Body()),
	}
}
