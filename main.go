package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pouya-eghbali/posh/pkg/lang/parser"
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

	if inputPath == "" {
		fmt.Println("No input file provided")
		os.Exit(1)
	}

	if outputPath == "" {
		fmt.Println("No output file provided")
		os.Exit(1)
	}

	err := parser.CompileFile(inputPath, outputPath, astOutput)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
