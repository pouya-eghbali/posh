package rules

import "github.com/pouya-eghbali/alien-go/pkg/lang/parser/types"

// TODO: Needs plug and unplug
type Assignment struct {
	types.BaseNode
	Identifier types.Node `json:"identifier"`
	Value      types.Node `json:"value"`
	ValueType  string     `json:"valueType"`
}

func MatchAssignment(nodes []types.Node, offset int) types.Result {
	start := offset

	if nodes[offset].GetType() != "IDENTIFIER" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	node := Assignment{
		BaseNode: types.BaseNode{
			Type: "ASSIGNMENT",
		},
		Identifier: nodes[offset],
	}

	offset++

	if nodes[offset].GetType() != "PUNCTUATOR" || nodes[offset].GetImage() != "=" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	offset++

	// We are looking for a simple expression
	if res := MatchExpr(nodes, offset); res.End > res.Start {
		node.Value = res.Node
		offset = res.End
	} else {
		return types.Result{FailedAt: res.FailedAt}
	}

	return types.Result{Node: &node, Start: start, End: offset}
}
