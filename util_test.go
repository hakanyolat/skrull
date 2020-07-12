package skrull

import (
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"testing"
)

func Test_statusIsMoved(t *testing.T) {
	movedStatuses := []int{
		fasthttp.StatusMovedPermanently,
		fasthttp.StatusFound,
		fasthttp.StatusSeeOther,
		fasthttp.StatusTemporaryRedirect,
		fasthttp.StatusPermanentRedirect,
	}

	for _, status := range movedStatuses {
		assert.Equal(t, true, statusIsMoved(status))
	}
}

func Test_createFullPath(t *testing.T) {
	assert.Equal(t, "test.com/test.html", createFullPath("test.com", "/test.html"))
}

func Test_createFullRelativePath(t *testing.T) {
	assert.Equal(t, "test.html?q=1", createFullRelativePath("test.html", "q=1"))
}