package types

import "go/ast"

type Param struct {
	Name string
	Type string
}

type Export struct {
	Type   string
	IsFunc bool
	Params []Param
}

type CompiledFile struct {
	FileName string
	Exports  map[string]Export
}

type PoshFile struct {
	Environment         *Environment
	Exports             map[string]Export
	TopLevelAssignments []ast.Spec
	CompiledFiles       map[string]CompiledFile
	StdImports          map[string]bool
	Source              string
	BaseDir             string
	OutputDir           string
	Package             string
}

func NewPoshFile(source string, basedir string, outputDir string, packageName string, compiledFiles map[string]CompiledFile) *PoshFile {
	return &PoshFile{
		Environment:         NewEnvironment(),
		TopLevelAssignments: []ast.Spec{},
		StdImports:          map[string]bool{},
		Exports:             map[string]Export{},
		CompiledFiles:       compiledFiles,
		Source:              source,
		BaseDir:             basedir,
		OutputDir:           outputDir,
		Package:             packageName,
	}
}

func NewCompiledFile(fileName string) CompiledFile {
	return CompiledFile{
		FileName: fileName,
		Exports:  map[string]Export{},
	}
}
