package skrull

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Miner(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("<div>Test</div>"))
	}))
	defer server.Close()

	app := NewApp(server.URL)
	app.Craft(func(parser *HtmlParser) {
		parser.Document.SetHtml("<div>Miner</div>")
	})
	app.Listen(8080)

	resp, _ := app.http.Test(httptest.NewRequest(http.MethodGet, "/", nil))
	body, _ := ioutil.ReadAll(resp.Body)
	mountedBody := "<div>Miner</div>"

	assert.Equal(t, mountedBody, string(body))
}