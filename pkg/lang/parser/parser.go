package parser

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"os/exec"
	"path"

	"github.com/pouya-eghbali/posh/pkg/lang/lexer"
	"github.com/pouya-eghbali/posh/pkg/lang/parser/rules"
	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
	"github.com/pouya-eghbali/posh/pkg/lang/parser/utils"
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

func ParseFile(path string) (error, string, types.Node, *types.Node) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		return fmt.Errorf("failed to read file: %v", err), "", &types.BaseNode{}, nil
	}

	code := string(bytes)
	parsed, failedAt := Parse(code)

	return nil, code, parsed, failedAt
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

func WriteResultToTempDir(node ast.Node) (string, error) {
	fset := token.NewFileSet()
	file := utils.GoAstToString(fset, node)

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

func CompileFile(filePath string, output string, printAst bool) error {
	err, code, parsed, failedAt := ParseFile(filePath)

	if err != nil {
		return err
	}

	if failedAt != nil {
		PrintError(code, failedAt)
		return nil
	}

	if printAst {
		PrintJSON(parsed)
	}

	goNode, poshFile := parsed.ToGoAstAndPoshFile("main")
	tempDir, err := WriteResultToTempDir(goNode)
	if err != nil {
		return err
	}

	compiled := map[string]bool{}
	toCompile := append([]string{}, poshFile.LocalImports...)
	mainDir := path.Dir(filePath)

	for len(toCompile) > 0 {
		imp := utils.Unquote(toCompile[0])
		impName := rules.ImportPathToImportName(toCompile[0])
		toCompile = toCompile[1:]

		if compiled[imp] {
			continue
		}

		impPath := path.Join(mainDir, imp)
		err, code, parsed, failedAt := ParseFile(impPath)

		if err != nil {
			return err
		}

		if failedAt != nil {
			PrintError(code, failedAt)
			return nil
		}

		goNode, poshFile := parsed.ToGoAstAndPoshFile(impName)
		outputDir := path.Join(tempDir, imp[:len(imp)-5])
		outputPath := path.Join(outputDir, "main.go")

		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create output dir: %v", err)
		}

		err = os.WriteFile(outputPath, []byte(utils.GoAstToString(token.NewFileSet(), goNode)), 0644)
		if err != nil {
			return fmt.Errorf("failed to write output file: %v", err)
		}

		toCompile = append(toCompile, poshFile.LocalImports...)
		compiled[imp] = true
	}

	err = CompileTempDir(tempDir, output)
	if err != nil {
		return err
	}

	return nil
}
