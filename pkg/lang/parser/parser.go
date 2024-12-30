package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"os"

	"github.com/pouya-eghbali/alien-go/pkg/lang/lexer"
	"github.com/pouya-eghbali/alien-go/pkg/lang/parser/rules"
	"github.com/pouya-eghbali/alien-go/pkg/lang/parser/types"
)

func Parse(code string) (types.Node, *types.Node) {
	tokens, err := lexer.Lex(code)

	if err != nil {
		fmt.Println(err)
		return &types.BaseNode{}, nil
	}

	res := rules.MatchAlien(tokens, 0)

	if res.FailedAt != nil {
		return &types.BaseNode{}, res.FailedAt
	}

	return res.Node, nil
}

func ParseFile(path string) (string, types.Node, *types.Node) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		fmt.Println(err)
		return "", &types.BaseNode{}, nil
	}

	code := string(bytes)
	parsed, failedAt := Parse(code)

	return code, parsed, failedAt
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

func GoAstToString(fset *token.FileSet, node ast.Node) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, node); err != nil {
		log.Fatalf("failed to print AST node: %v", err)
	}
	return buf.String()
}
