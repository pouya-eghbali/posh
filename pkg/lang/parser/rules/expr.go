package rules

import (
	"github.com/pouya-eghbali/alien-go/pkg/lang/parser/types"
)

type SimpleExpression struct {
	types.BaseNode
	Value types.Node `json:"value"`
}

func MatchExpr(nodes []types.Node, offset int) types.Result {
	// We are looking for Pipe or simple expression

	// look for PIPE
	if res := MatchPipe(nodes, offset); res.End > res.Start {
		return res
	}

	// look for simple expression
	if res := MatchSimpleExpression(nodes, offset); res.End > res.Start {
		return res
	} else {
		return types.Result{FailedAt: res.FailedAt}
	}

}

func MatchSimpleExpression(nodes []types.Node, offset int) types.Result {
	start := offset

	// We are looking for one of the following:
	// - ARITHMETIC
	// - CALL
	// - NUMERIC
	// - STRING
	// - BOOLEAN

	// try to match ARITHMETIC
	if res := MatchArithmetic(nodes, offset); res.End > res.Start {
		return res
	}

	// try to match CALL
	if res := MatchFunctionCall(nodes, offset); res.End > res.Start {
		return res
	}

	// try to match NUMERIC
	if res := MatchNumeric(nodes, offset); res.End > res.Start {
		return res
	}

	// try to match the rest
	if nodes[offset].GetType() != "STRING" && nodes[offset].GetType() != "BOOLEAN" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	node := SimpleExpression{
		BaseNode: types.BaseNode{
			Type: "SIMPLE_EXPRESSION",
		},
		Value: nodes[offset],
	}

	return types.Result{Node: &node, Start: start, End: offset + 1}
}
