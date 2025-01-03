package rules

import (
	"go/ast"
	"go/token"

	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

type Logical struct {
	types.BaseNode
	Lhs types.Node `json:"lhs"`
	Op  types.Node `json:"op"`
	Rhs types.Node `json:"rhs"`
}

var logicalTokenMap = map[string]token.Token{
	"and": token.LAND,
	"or":  token.LOR,
}

func isLogicalOperator(value string) bool {
	return value == "and" || value == "or"
}

func (n *Logical) ToGoAst() ast.Node {
	return &ast.BinaryExpr{
		X:  n.Lhs.ToGoAst().(ast.Expr),
		Op: logicalTokenMap[n.Op.GetImage()],
		Y:  n.Rhs.ToGoAst().(ast.Expr),
	}
}

func MatchLogical(nodes []types.Node, offset int) types.Result {
	start := offset

	// We should recursively match the following:
	// BOOLEAN (LOGICAL_OPERATOR BOOLEAN)*

	logical := Logical{
		BaseNode: types.BaseNode{
			Type: "LOGICAL",
		},
	}

	if res := MatchBoolean(nodes, offset); res.End <= res.Start {
		return types.Result{FailedAt: &nodes[offset]}
	} else {
		offset = res.End
		logical.Lhs = res.Node
	}

	// match LOGICAL_OPERATOR
	if nodes[offset].GetType() != "KEYWORD" || !isLogicalOperator(nodes[offset].GetImage()) {
		return types.Result{FailedAt: &nodes[offset]}
	}
	logical.Op = nodes[offset]
	offset++

	// match BOOLEAN
	if res := MatchBoolean(nodes, offset); res.End > res.Start {
		offset = res.End
		logical.Rhs = res.Node
	} else if res := MatchLogical(nodes, offset); res.End > res.Start {
		offset = res.End
		logical.Rhs = res.Node
	} else {
		return types.Result{FailedAt: &nodes[offset]}
	}

	return types.Result{Node: &logical, Start: start, End: offset + 1}
}
