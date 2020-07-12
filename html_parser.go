package skrull

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"log"
)

type HtmlParser struct {
	Document *goquery.Document
	zap.Logger
}

// NewHtmlParser ...
func NewHtmlParser(res []byte) *HtmlParser {
	reader := bytes.NewReader(res)
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err)
	}
	return &HtmlParser{Document: document}
}

// Replace ...
func (h HtmlParser) Replace(fragment RenderedFragment) {
	if fragment.Result != "" {
		h.Document.Find(fragment.Selector).ReplaceWithHtml(fragment.Result)
	}
}
