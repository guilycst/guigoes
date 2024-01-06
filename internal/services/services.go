package services

import (
	"bytes"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/anchor"
)

func markdownToHTML(md []byte) []byte {
	parser := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // read note
		),
		goldmark.WithExtensions(
			extension.Table,
			extension.Strikethrough,
			extension.Linkify,
			extension.TaskList,
			extension.DefinitionList,
			extension.Footnote,
			extension.Typographer,
			highlighting.NewHighlighting(
				highlighting.WithStyle("xcode-dark"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
			&anchor.Extender{
				Texter:   anchor.Text("🔗"),
				Position: anchor.Before,
			},
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	// create markdown parser with extensions
	var buf bytes.Buffer
	if err := parser.Convert(md, &buf); err != nil {
		panic(err)
	}

	return buf.Bytes()
}
