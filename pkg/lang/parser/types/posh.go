package types

import "go/ast"

type PoshFile struct {
	Environment         *Environment
	TopLevelAssignments []ast.Spec
	LocalImports        []string
	StdImports          map[string]bool
}

func NewPoshFile() *PoshFile {
	return &PoshFile{
		Environment:         NewEnvironment(),
		TopLevelAssignments: []ast.Spec{},
		LocalImports:        []string{},
		StdImports:          map[string]bool{},
	}
}
