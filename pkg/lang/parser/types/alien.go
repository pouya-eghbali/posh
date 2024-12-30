package types

import "go/ast"

type PoshFile struct {
	Environment         *Environment
	TopLevelAssignments []ast.Spec
}

func NewPoshFile() *PoshFile {
	return &PoshFile{
		Environment:         NewEnvironment(),
		TopLevelAssignments: []ast.Spec{},
	}
}
