package types

import "go/ast"

type Pos struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Node interface {
	GetPos() *Pos
	GetImage() string
	GetType() string
	ToGoAst() ast.Node
	CompileToGo(posh *PoshFile) error
	ToGoStatementAst() ast.Stmt
	StaticAnalysis(*PoshFile)
	Plug(*Environment)
	UnPlug(*Environment)
}

type TopLevelMatcher = func(nodes []Node, offset int) Result

type Result struct {
	Node     Node
	FailedAt *Node
	Start    int
	End      int
}
