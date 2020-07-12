package skrull

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApp_RequestMiddleware(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("Body"))
	}))
	defer server.Close()

	app := NewApp(server.URL)
	app.RequestMiddleware(func(ctx *MiddlewareContext) {
		ctx.Ctx.Status(404)
	})
	app.Listen(8080)

	resp, _ := app.http.Test(httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, 404, resp.StatusCode)
}

func TestApp_ResponseMiddleware(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("Body"))
	}))
	defer server.Close()

	app := NewApp(server.URL)
	app.ResponseMiddleware(func(ctx *MiddlewareContext) {
		ctx.Res.SetBodyString("Mounted Body")
	})
	app.Listen(8080)

	resp, _ := app.http.Test(httptest.NewRequest(http.MethodGet, "/", nil))

	body, _ := ioutil.ReadAll(resp.Body)
	mountedBody := "<html><head></head><body>Mounted Body</body></html>"
	assert.Equal(t, mountedBody, string(body))
}
