// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.648
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func CarrousselComponent() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"carousel w-full\"><div id=\"slide1\" class=\"carousel-item relative w-full\"><img src=\"https://daisyui.com/images/stock/photo-1625726411847-8cbb60cc71e6.jpg\" class=\"w-full\"><div class=\"absolute flex justify-between transform -translate-y-1/2 left-5 right-5 top-1/2\"><a href=\"#slide4\" class=\"btn btn-circle\">❮</a> <a href=\"#slide2\" class=\"btn btn-circle\">❯</a></div></div><div id=\"slide2\" class=\"carousel-item relative w-full\"><img src=\"https://daisyui.com/images/stock/photo-1609621838510-5ad474b7d25d.jpg\" class=\"w-full\"><div class=\"absolute flex justify-between transform -translate-y-1/2 left-5 right-5 top-1/2\"><a href=\"#slide1\" class=\"btn btn-circle\">❮</a> <a href=\"#slide3\" class=\"btn btn-circle\">❯</a></div></div><div id=\"slide3\" class=\"carousel-item relative w-full\"><img src=\"https://daisyui.com/images/stock/photo-1414694762283-acccc27bca85.jpg\" class=\"w-full\"><div class=\"absolute flex justify-between transform -translate-y-1/2 left-5 right-5 top-1/2\"><a href=\"#slide2\" class=\"btn btn-circle\">❮</a> <a href=\"#slide4\" class=\"btn btn-circle\">❯</a></div></div><div id=\"slide4\" class=\"carousel-item relative w-full\"><img src=\"https://daisyui.com/images/stock/photo-1665553365602-b2fb8e5d1707.jpg\" class=\"w-full\"><div class=\"absolute flex justify-between transform -translate-y-1/2 left-5 right-5 top-1/2\"><a href=\"#slide3\" class=\"btn btn-circle\">❮</a> <a href=\"#slide1\" class=\"btn btn-circle\">❯</a></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
