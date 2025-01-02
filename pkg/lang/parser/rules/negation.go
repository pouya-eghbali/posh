package rules

import (
	"go/ast"
	"go/token"

	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

type Negation struct {
	types.BaseNode
	Value types.Node `json:"value"`
}

func (n *Negation) ToGoAst() ast.Node {
	return &ast.UnaryExpr{
		Op: token.NOT,
		X:  n.Value.ToGoAst().(ast.Expr),
	}
}

func MatchNegation(nodes []types.Node, offset int) types.Result {
	start := offset

	// We are looking for one of the following:
	// NOT BOOLEAN

	// try to match NOT
	if nodes[offset].GetType() != "KEYWORD" || nodes[offset].GetImage() != "not" {
		return types.Result{FailedAt: &nodes[offset]}
	}
	offset++

	// try to match BOOLEAN
	var res types.Result
	if res = MatchBoolean(nodes, offset); res.End <= res.Start {
		return types.Result{FailedAt: &nodes[offset]}
	}
	offset = res.End

	node := Negation{
		BaseNode: types.BaseNode{
			Type: "NEGATION",
		},
		Value: res.Node,
	}

	return types.Result{Node: &node, Start: start, End: offset + 1}
}
