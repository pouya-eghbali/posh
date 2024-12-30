package main

import (
	"flag"

	"github.com/pouya-eghbali/alien-go/pkg/lang/parser"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "file", "", "Path to the file to parse")

	var jsonOutput bool
	flag.BoolVar(&jsonOutput, "json", false, "Output the result in JSON format")
	flag.Parse()

	code, node, failedAt := parser.ParseFile(filePath)
	if failedAt != nil {
		parser.PrintError(code, failedAt)
		return
	}

	if jsonOutput {
		parser.PrintJSON(node)
	} else {
		parser.Print(node)
	}

}
