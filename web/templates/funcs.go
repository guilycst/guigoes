package templates

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/golang-module/carbon/v2"
)

func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}

func LocalizeTime(dt carbon.DateTime, lang string) string {
	return dt.SetLocale(lang).Format("M, d, Y")
}
