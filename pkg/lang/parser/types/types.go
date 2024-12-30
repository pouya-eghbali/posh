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
	ToGoStatementAst() ast.Stmt
	CollectTopLevelAssignments(*PoshFile)
	Plug(*Environment)
	UnPlug(*Environment)
}

type Result struct {
	Node     Node
	FailedAt *Node
	Start    int
	End      int
}
