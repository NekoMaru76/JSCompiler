package main

import (
	"fmt"
	"os"

	js_to_go "github.com/NekoMaru76/JSCompiler/collections/go"
	parser "github.com/NekoMaru76/JSCompiler/grammars/lib"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

func main() {
	input, _ := antlr.NewFileStream(os.Args[2])
	lexer := parser.NewJavaScriptLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewJavaScriptParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true
	visitor := js_to_go.Visitor{}
	prog := p.Program().(*parser.ProgramContext)
	ret := visitor.VisitProgram(prog)

	fmt.Println(ret.ToFileString(os.Args[1]))
	err := os.WriteFile(os.Args[3], []byte(ret.ToFileString(os.Args[1])), 0644)

	if err == nil {
		return
	}

	fmt.Println(err)
}
