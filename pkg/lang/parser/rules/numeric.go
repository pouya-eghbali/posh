package rules

import "github.com/pouya-eghbali/alien-go/pkg/lang/parser/types"

type Numeric struct {
	types.BaseNode
	Value types.Node `json:"value"`
}

func MatchNumeric(nodes []types.Node, offset int) types.Result {
	start := offset

	// We are looking for one of the following:
	// - WRAPPED_ARITHMETIC
	// - DOT_NOTATION
	// - INTEGER
	// - FLOAT
	// - IDENTIFIER

	// try to match WRAPPED_ARITHMETIC
	if res := MatchWrappedArithmetic(nodes, offset); res.End > res.Start {
		return res
	}

	// try to match DOT_NOTATION
	if res := MatchDotNotation(nodes, offset); res.End > res.Start {
		return res
	}

	if nodes[offset].GetType() != "INTEGER" &&
		nodes[offset].GetType() != "FLOAT" &&
		nodes[offset].GetType() != "IDENTIFIER" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	node := Numeric{
		BaseNode: types.BaseNode{
			Type: "NUMERIC",
		},
		Value: nodes[offset],
	}

	return types.Result{Node: &node, Start: start, End: offset + 1}
}
