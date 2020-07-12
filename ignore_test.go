package skrull

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApp_Ignore(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("<div id='test-fragment'></div>"))
	}))
	defer server.Close()

	testFragment := Fragment{
		Name:     "Test Fragment",
		Url:      "http://test-fragment.example.com",
		Selector: "#test-fragment",
	}

	app := NewApp(server.URL)
	app.Replace(testFragment)
	app.PassAll("/")
	app.Listen(8080)

	resp, _ := app.http.Test(httptest.NewRequest(http.MethodGet, "/", nil))
	body, _ := ioutil.ReadAll(resp.Body)
	mountedBody := "<div id='test-fragment'></div>"

	assert.Equal(t, mountedBody, string(body))
}