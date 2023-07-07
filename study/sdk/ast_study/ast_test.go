package ast_study

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestAST(t *testing.T) {
	src := `
package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	ast.Print(fset, file)
}
