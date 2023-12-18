package services

import (
	"bytes"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

func markdownToHTML(md []byte) []byte {

	parser := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("friendly"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
	)

	// create markdown parser with extensions
	var buf bytes.Buffer
	if err := parser.Convert(md, &buf); err != nil {
		panic(err)
	}

	return buf.Bytes()
}
