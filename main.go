package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/k0kubun/pp"
	"github.com/tzmfreedom/land/ast"
	"github.com/tzmfreedom/land/parser"
)

func main() {
	file := os.Args[1]
	node, err := ParseFile(file)
	if err != nil {
		panic(err)
	}
	resolver := NewImportTypeResolver()
	_, err = resolver.Resolve(node)
	if err != nil {
		panic(err)
	}
	for importClass, _ := range resolver.importClasses {
		fmt.Printf("import %s;\n", importClass)
	}
	fmt.Println()
	fmt.Println(Generate(node))
}

func ParseFile(f string) (ast.Node, error) {
	bytes, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	input := antlr.NewInputStream(string(bytes))
	return parse(input, f), nil
}

func ParseString(src string) (ast.Node, error) {
	input := antlr.NewInputStream(src)
	return parse(input, "<string>"), nil
}

func parse(input antlr.CharStream, src string) ast.Node {
	lexer := parser.NewapexLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewapexParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true
	tree := p.CompilationUnit()
	t := tree.Accept(&ast.Builder{
		Source: src,
	})
	return t.(ast.Node)
	//return nil
}

func debug(args ...interface{}) {
	pp.Println(args...)
}
