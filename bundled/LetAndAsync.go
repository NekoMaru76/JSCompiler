package bundled

import js_to_go "github.com/NekoMaru76/JSCompiler/collections/go"

func main() {
	scopeCtx := js_to_go.Context{}
	funcCtx := js_to_go.Context{}

	scopeCtx.DefineVariable("foo", &js_to_go.Number{Value: 123}, false)
	funcCtx.DefineVariable("let", &js_to_go.Number{Value: 123}, false)
	funcCtx.DefineVariable("async", &js_to_go.Number{Value: 456}, false)
	&js_to_go.String{Value: "use strict"}
	funcCtx.DefineVariable("let", &js_to_go.Number{Value: 123}, false)
}
