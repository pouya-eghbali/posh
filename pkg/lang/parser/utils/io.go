package utils

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path"

	"github.com/pouya-eghbali/posh/pkg/lang/lexer"
	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

func Parse(code string, topLevel types.TopLevelMatcher) (types.Node, *types.Node) {
	tokens, err := lexer.Lex(code)

	if err != nil {
		fmt.Println(err)
		return &types.BaseNode{}, nil
	}

	res := topLevel(tokens, 0)
	if res.FailedAt != nil {
		return &types.BaseNode{}, res.FailedAt
	}

	return res.Node, nil
}

func ParseFile(path string, topLevel types.TopLevelMatcher) (error, string, types.Node, *types.Node) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		return fmt.Errorf("failed to read file: %v", err), "", &types.BaseNode{}, nil
	}

	code := string(bytes)
	parsed, failedAt := Parse(code, topLevel)

	return nil, code, parsed, failedAt
}

func WritePoshFile(node ast.Node, posh *types.PoshFile) error {
	outputDir := posh.OutputDir

	if posh.Package != "main" {
		outputDir = path.Join(outputDir, path.Dir(posh.Source), posh.Package)
	}

	outputPath := path.Join(outputDir, "main.go")

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output dir: %v", err)
	}

	err = os.WriteFile(outputPath, []byte(GoAstToString(token.NewFileSet(), node)), 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	return nil
}

func CompilePoshFile(posh *types.PoshFile, topLevel types.TopLevelMatcher) error {
	filePath := path.Join(posh.BaseDir, posh.Source)
	err, code, parsed, failedAt := ParseFile(filePath, topLevel)

	if err != nil {
		return err
	}

	if failedAt != nil {
		PrintError(code, failedAt)
		return nil
	}

	return parsed.CompileToGo(posh)
}
