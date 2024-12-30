package rules

import (
	types "github.com/pouya-eghbali/alien-go/pkg/lang/parser/types"
)

type DotNotation struct {
	types.BaseNode
	Accessors []types.Node `json:"accessors"`
}

func MatchDotNotation(nodes []types.Node, offset int) types.Result {
	start := offset

	if nodes[offset].GetType() != "IDENTIFIER" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	// if there are no dots, this is not a dot notation
	if offset+1 >= len(nodes) || nodes[offset+1].GetType() != "PUNCTUATOR" || nodes[offset+1].GetImage() != "." {
		return types.Result{FailedAt: &nodes[offset+1]}
	}

	node := DotNotation{
		BaseNode: types.BaseNode{
			Type: "DOT_NOTATION",
		},
	}

	offset++

	// recursively match the dots and identifiers
	for {
		if nodes[offset].GetType() != "PUNCTUATOR" || nodes[offset].GetImage() != "." {
			break
		}

		offset++

		if nodes[offset].GetType() != "IDENTIFIER" {
			return types.Result{FailedAt: &nodes[offset]}
		}

		node.Accessors = append(node.Accessors, nodes[offset])

		offset++
	}

	return types.Result{Node: &node, Start: start, End: offset}
}
