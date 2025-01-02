package rules

import (
	"go/ast"

	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

type Boolean struct {
	types.BaseNode
	Value types.Node `json:"value"`
}

func (n *Boolean) ToGoAst() ast.Node {
	return ast.NewIdent(n.Value.GetImage())
}

func MatchBoolean(nodes []types.Node, offset int) types.Result {
	start := offset

	// We are looking for one of the following:
	// - COMPARISON
	// - NUMERIC
	// - TRUE|FALSE

	// try to match COMPARISON
	if res := MatchComparison(nodes, offset); res.End > res.Start {
		return res
	}

	// try to match NUMERIC
	if res := MatchNumeric(nodes, offset); res.End > res.Start {
		return res
	}

	if nodes[offset].GetImage() != "true" && nodes[offset].GetImage() != "false" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	node := Boolean{
		BaseNode: types.BaseNode{
			Type: "BOOLEAN",
		},
		Value: nodes[offset],
	}

	return types.Result{Node: &node, Start: start, End: offset + 1}
}
