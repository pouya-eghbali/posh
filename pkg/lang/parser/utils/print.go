package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"log"

	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

func GoAstToString(fset *token.FileSet, node ast.Node) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, node); err != nil {
		log.Fatalf("failed to print AST node: %v", err)
	}
	return buf.String()
}

func PrintJSON(node types.Node) {
	res, err := json.MarshalIndent(node, "", "  ")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(res))
}

func Print(node types.Node) {
	fmt.Println(node)
}
