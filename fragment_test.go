package skrull

import (
	"github.com/gofiber/fiber"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFragment_Available(t *testing.T) {
	parser  := NewHtmlParser([]byte(`<html><head></head><body><div id="test-fragment"></div></html>`))
	testFragment := Fragment{
		Name:     "Test Fragment",
		Url:      "http://test-fragment.example.com",
		Selector: "#test-fragment",
	}
	assert.True(t, testFragment.Available(parser))
}

func TestFragment_Select(t *testing.T) {
	parser  := NewHtmlParser([]byte(`<html><head></head><body><div id="test-fragment"></div></html>`))
	testFragment := Fragment{
		Name:     "Test Fragment",
		Url:      "http://test-fragment.example.com",
		Selector: "#test-fragment",
	}

	selection := testFragment.Select(parser)
	assert.True(t, selection.IsSelection(parser.Document.Find(testFragment.Selector)))
}

func TestFragment_Render(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("<div id='test-fragment-replaced'></div>"))
	}))
	defer server.Close()

	fragments := []Fragment{
		{
			Name: "Test",
			Url: server.URL,
			Selector: "#test-fragment-holder",
		},
	}

	parser  := NewHtmlParser([]byte(`<html><head></head><body><div id="test-fragment-holder"></div></html>`))

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)
	mCtx := MiddlewareContext{
		Ctx:   ctx,
	}

	ch := make(chan RenderedFragment)
	for _, fragment := range fragments {
		if fragment.Available(parser) {
			go fragment.Render(mCtx, ch, nil)
		}
	}

	for range fragments {
		res := <-ch
		parser.Replace(res)
	}

	replacedElement := parser.Document.Find("#test-fragment-replaced").Get(0)
	assert.True(t, parser.Document.Contains(replacedElement))
}