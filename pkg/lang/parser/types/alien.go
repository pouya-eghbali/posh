package types

import "go/ast"

type AlienFile struct {
	Environment         *Environment
	TopLevelAssignments []ast.Spec
}

func NewAlienFile() *AlienFile {
	return &AlienFile{
		Environment:         NewEnvironment(),
		TopLevelAssignments: []ast.Spec{},
	}
}
