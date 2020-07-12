package skrull

import (
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"testing"
)

func TestApp_Craft(t *testing.T) {
	app := NewApp("localhost")
	assert.Equal(t, 0, len(app.miners))

	app.Craft(func(parser *HtmlParser) {})
	assert.Equal(t, 1, len(app.miners))
}

func TestApp_Pass(t *testing.T) {
	app := NewApp("localhost")
	assert.Equal(t, 0, len(app.ignores))

	app.Pass(fasthttp.MethodGet, "/test.html")
	assert.Equal(t, 1, len(app.ignores))
}

func TestApp_PassAll(t *testing.T) {
	app := NewApp("localhost")
	assert.Equal(t, 0, len(app.ignores))

	app.PassAll("/test.html")
	assert.Equal(t, len(httpMethods), len(app.ignores))
}

func TestApp_Listen(t *testing.T) {
	app := NewApp("localhost")
	assert.Nil(t, app.address)

	app.Listen(8080)
	assert.Equal(t, app.address, 8080)
	assert.NotNil(t, app.http)
}

func TestApp_Replace(t *testing.T) {
	app := NewApp("localhost")
	assert.Equal(t, 0, len(app.fragments))

	app.Replace(Fragment{
		Name:     "Test Fragment",
		Url:      "http://example.com",
		Selector: "#test-fragment",
	})

	assert.Equal(t, 1, len(app.fragments))
}

func TestApp_SyncCookie(t *testing.T) {
	app := NewApp("localhost")
	assert.Equal(t, 0, len(app.cookieTable))

	app.SyncCookie("test_cookie")
	assert.Equal(t, 1, len(app.cookieTable))
}