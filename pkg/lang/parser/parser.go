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
	"os/exec"
	"path"

	"github.com/pouya-eghbali/posh/pkg/lang/lexer"
	"github.com/pouya-eghbali/posh/pkg/lang/parser/rules"
	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

func Parse(code string) (types.Node, *types.Node) {
	tokens, err := lexer.Lex(code)

	if err != nil {
		fmt.Println(err)
		return &types.BaseNode{}, nil
	}

	res := rules.MatchPosh(tokens, 0)

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

func WriteResultToTempDir(node types.Node) (string, error) {
	fset := token.NewFileSet()
	file := GoAstToString(fset, node.ToGoAst())

	tempDir, err := os.MkdirTemp("", "posh-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %v", err)
	}

	tempFile := path.Join(tempDir, "main.go")
	err = os.WriteFile(tempFile, []byte(file), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write temp file: %v", err)
	}

	return tempDir, nil
}

func CompileTempDir(tempDir string, output string) error {
	// Run go mod init
	cmd := exec.Command("go", "mod", "init", "main")
	cmd.Dir = tempDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run go mod init: %v, output: %s", err, string(output))
	}

	// Run go mod tidy
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = tempDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %v, output: %s", err, string(output))
	}

	// Run go build
	cmd = exec.Command("go", "build", "-ldflags", "-s -w", "-o", "main")
	cmd.Dir = tempDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to compile temp dir: %v, output: %s", err, string(output))
	}

	// Move the compiled binary to the output path
	err := os.Rename(path.Join(tempDir, "main"), output)
	if err != nil {
		return fmt.Errorf("failed to move compiled binary: %v", err)
	}

	// Remove the temp dir
	err = os.RemoveAll(tempDir)
	if err != nil {
		return fmt.Errorf("failed to remove temp dir: %v", err)
	}

	return nil
}

func CompileFile(path string, output string, printAst bool) error {
	code, parsed, failedAt := ParseFile(path)

	if failedAt != nil {
		PrintError(code, failedAt)
		return nil
	}

	if printAst {
		PrintJSON(parsed)
	}

	tempDir, err := WriteResultToTempDir(parsed)
	if err != nil {
		return err
	}

	err = CompileTempDir(tempDir, output)
	if err != nil {
		return err
	}

	return nil
}
