// Code generated by templ - DO NOT EDIT.

// templ: version: 0.2.476
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/guilycst/guigoes/web/templates/state"

func Base(bs state.BaseState) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"pt-BR\"><head><title>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string = bs.Title
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</title><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><meta name=\"description\" content=\"A developer blog sharing insights, tutorials, and tips on programming and web development.\"><meta name=\"author\" content=\"Guilherme de Castro\"><script src=\"https://unpkg.com/htmx.org@1.9.9\" integrity=\"sha384-QFjmbokDn2DjBjq+fM+8LUIVrAgqcNW2s0PjAxHETgRn9l4fvX31ZxDxvwQnyMOX\" crossorigin=\"anonymous\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var3 := ``
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var3)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</script><link rel=\"stylesheet\" href=\"/output.css\" defer><link rel=\"apple-touch-icon\" sizes=\"180x180\" href=\"/apple-touch-icon.png\"><link rel=\"icon\" type=\"image/png\" sizes=\"32x32\" href=\"/favicon-32x32.png\"><link rel=\"icon\" type=\"image/png\" sizes=\"16x16\" href=\"/favicon-16x16.png\"><link rel=\"manifest\" href=\"/site.webmanifest\"><link rel=\"preconnect\" href=\"https://fonts.googleapis.com\" defer><link rel=\"preconnect\" href=\"https://fonts.gstatic.com\" crossorigin defer></head><body class=\"bg-zinc-900 text-zinc-300 text-xl\"><nav class=\"px-2 flex flex-row items-center bg-zinc-950\"><a href=\"/\" class=\"text-white text-2xl font-bold\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var4 := `GG`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var4)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a><div class=\"flex items-center justify-center h-16 m-auto md:max-w-prose\"><div class=\"flex items-center\"><div><div class=\"flex flex-row items-center space-x-4\"><form autocomplete=\"off\"><input autocomplete=\"off\" name=\"hidden\" type=\"text\" class=\"hidden\"><div class=\"relative\"><svg class=\"absolute top-0 left-2 bottom-0 mt-auto mr-0 mb-auto ml-0\" width=\"24px\" height=\"24px\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\" color=\"#c2c1c1\"><path d=\"M17 17L21 21\" stroke=\"#c2c1c1\" stroke-width=\"1.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\"></path><path d=\"M3 11C3 15.4183 6.58172 19 11 19C13.213 19 15.2161 18.1015 16.6644 16.6493C18.1077 15.2022 19 13.2053 19 11C19 6.58172 15.4183 3 11 3C6.58172 3 3 6.58172 3 11Z\" stroke=\"#c2c1c1\" stroke-width=\"1.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\"></path></svg> <input autocomplete=\"off\" class=\"pl-9 md:block bg-zinc-900 text-sm font-medium text-zinc-300 rounded-full\" type=\"text\" name=\"search\" id=\"search\" placeholder=\"Search\" hx-trigger=\"keyup changed delay:500ms\" hx-post=\"/search\" hx-target=\"#main\" hx-swap=\"innerHTML\"></div></form><a href=\"/\" class=\"hidden md:block text-zinc-300 hover:bg-zinc-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var5 := `Home`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var5)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a> <a hx-get=\"/about?fragment=1\" hx-target=\"#main\" class=\"hidden md:block text-zinc-300 hover:bg-zinc-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var6 := `About`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var6)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a></div></div></div></div><div class=\"group relative md:hidden hover:bg-zinc-700 hover:text-white rounded-t-md px-2 justify-self-end\"><svg width=\"40px\" height=\"40px\" stroke-width=\"1.5\" viewBox=\"0 0 24 24\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\" color=\"#deded8\"><path d=\"M3 5H21\" stroke=\"#deded8\" stroke-width=\"1.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\"></path><path d=\"M3 12H21\" stroke=\"#deded8\" stroke-width=\"1.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\"></path><path d=\"M3 19H21\" stroke=\"#deded8\" stroke-width=\"1.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\"></path></svg><div class=\"hidden group-hover:flex flex-col items-center absolute w-[100px] mw-[100px] top-[39px] -left-11 text-ellipsis bg-zinc-700 rounded-b-md rounded-tl-md\"><div class=\"hidden group-hover:block bg-zinc-700 h-7\"></div><a href=\"/\" class=\"mw-full text-zinc-300 hover:bg-zinc-700 hover:text-white rounded-md text-sm pb-2 font-medium\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var7 := `Home`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var7)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a> <a hx-get=\"/about?fragment=1\" class=\"mt-2 text-zinc-300 hover:bg-zinc-700 hover:text-white rounded-md text-sm pb-2 font-medium\" hx-target=\"#main\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var8 := `About`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var8)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a></div></div></nav><main id=\"main\" class=\"flex flex-col gap-4 place-content-center mb-10\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = bs.Body.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</main><link href=\"https://fonts.googleapis.com/css2?family=Montserrat:ital,wght@0,300;1,300&amp;display=swap\" rel=\"stylesheet\" defer></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
