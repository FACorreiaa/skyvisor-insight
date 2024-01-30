// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func TableComponent() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"relative overflow-x-auto\"><table class=\"w-full text-sm text-left text-gray-500 rtl:text-right dark:text-gray-400\"><thead class=\"text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400\"><tr><th scope=\"col\" class=\"px-6 py-3\">Product name</th><th scope=\"col\" class=\"px-6 py-3\">Color</th><th scope=\"col\" class=\"px-6 py-3\">Category</th><th scope=\"col\" class=\"px-6 py-3\">Price</th></tr></thead> <tbody><tr class=\"bg-white border-b dark:bg-gray-800 dark:border-gray-700\"><th scope=\"row\" class=\"px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white\">Apple MacBook Pro 17\"</th><td class=\"px-6 py-4\">Silver</td><td class=\"px-6 py-4\">Laptop</td><td class=\"px-6 py-4\">$2999</td></tr><tr class=\"bg-white border-b dark:bg-gray-800 dark:border-gray-700\"><th scope=\"row\" class=\"px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white\">Microsoft Surface Pro</th><td class=\"px-6 py-4\">White</td><td class=\"px-6 py-4\">Laptop PC</td><td class=\"px-6 py-4\">$1999</td></tr><tr class=\"bg-white dark:bg-gray-800\"><th scope=\"row\" class=\"px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white\">Magic Mouse 2</th><td class=\"px-6 py-4\">Black</td><td class=\"px-6 py-4\">Accessories</td><td class=\"px-6 py-4\">$99</td></tr></tbody></table></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
