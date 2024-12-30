package rules

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
)

func GoAstToString(fset *token.FileSet, node ast.Node) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, node); err != nil {
		log.Fatalf("failed to print AST node: %v", err)
	}
	return buf.String()
}
