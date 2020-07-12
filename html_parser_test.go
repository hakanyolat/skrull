package skrull

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHtmlParser_Replace(t *testing.T) {
	parser  := NewHtmlParser([]byte(`<html><head></head><body><div id="test"></div></html>`))
	fragment := RenderedFragment{
		Fragment: Fragment{
			Name: "Test Fragment",
			Url: "http://test-fragment.example.com",
			Selector: "#test",
		},
		Result:   "<div id='test-rendered'></div>",
	}
	parser.Replace(fragment)

	replacedElement := parser.Document.Find("#test-rendered").Get(0)
	assert.True(t, parser.Document.Contains(replacedElement))
}

