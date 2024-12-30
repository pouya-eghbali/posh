package main

import (
	"flag"
	"fmt"
	"go/token"
	"os"

	"github.com/pouya-eghbali/alien-go/pkg/lang/parser"
)

var usage = `Usage: alc [options]

Options:
	-i, -input string
		Path to the file to parse

	-o, -output string
		Path to the file to output

	-ast
		Output the ast in JSON format

Example:
	alc -i file.alc

`

func main() {
	var inputPath string
	flag.StringVar(&inputPath, "input", "", "Path to the file to parse")
	flag.StringVar(&inputPath, "i", "", "Path to the file to parse")

	var outputPath string
	flag.StringVar(&outputPath, "output", "", "Path to the file to output")
	flag.StringVar(&outputPath, "o", "", "Path to the file to output")

	var astOutput bool
	flag.BoolVar(&astOutput, "ast", false, "Output the ast in JSON format")

	flag.Usage = func() {
		fmt.Println(usage)
	}

	flag.Parse()

	code, node, failedAt := parser.ParseFile(inputPath)
	if failedAt != nil {
		parser.PrintError(code, failedAt)
		return
	}

	if astOutput {
		parser.PrintJSON(node)
	}

	if outputPath != "" {
		goAst := node.ToGoAst()
		out := parser.GoAstToString(token.NewFileSet(), goAst)
		err := os.WriteFile(outputPath, []byte(out), 0644)
		if err != nil {
			fmt.Println("Failed to write to the output file")
		}
	}

}
