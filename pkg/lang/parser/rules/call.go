package rules

import "github.com/pouya-eghbali/alien-go/pkg/lang/parser/types"

type Flag struct {
	types.BaseNode
	Identifier types.Node `json:"identifier"`
	DashCount  int        `json:"dashCount"`
}

type FunctionCall struct {
	types.BaseNode
	Callable types.Node   `json:"callable"`
	Args     []types.Node `json:"args"`
}

func MatchFlag(nodes []types.Node, offset int) types.Result {
	start := offset

	// a flag is an identifier prefixed with one or two dashes
	// --identifier or -identifier

	// check for the first dash
	if nodes[offset].GetType() != "PUNCTUATOR" || nodes[offset].GetImage() != "-" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	offset++

	node := Flag{
		BaseNode: types.BaseNode{
			Type: "FLAG",
		},
		DashCount: 1,
	}

	// check for the second dash
	if nodes[offset].GetType() == "PUNCTUATOR" && nodes[offset].GetImage() == "-" {
		node.DashCount = 2
		offset++
	}

	if nodes[offset].GetType() != "IDENTIFIER" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	node.Identifier = nodes[offset]
	return types.Result{Node: &node, Start: start, End: offset + 1}
}

func MatchFunctionCall(nodes []types.Node, offset int) types.Result {
	start := offset

	// We are looking for the following as callables:
	// - DOT_NOTATION
	// - IDENTIFIER

	// try to match DOT_NOTATION
	var callable types.Node
	if res := MatchDotNotation(nodes, offset); res.End > res.Start {
		callable = res.Node
		offset = res.End
	} else if nodes[offset].GetType() == "IDENTIFIER" {
		callable = nodes[offset]
		offset++
	} else {
		return types.Result{FailedAt: res.FailedAt}
	}

	node := FunctionCall{
		BaseNode: types.BaseNode{
			Type: "FUNCTION_CALL",
		},
		Callable: callable,
	}

	if nodes[offset].GetType() != "PUNCTUATOR" || nodes[offset].GetImage() != "(" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	offset++

	for {
		if nodes[offset].GetType() == "PUNCTUATOR" && nodes[offset].GetImage() == ")" {
			break
		}

		// Look for simple expressions or flags
		if res := MatchSimpleExpression(nodes, offset); res.End > res.Start {
			node.Args = append(node.Args, res.Node)
			offset = res.End
		} else if res := MatchFlag(nodes, offset); res.End > res.Start {
			node.Args = append(node.Args, res.Node)
			offset = res.End
		} else {
			return types.Result{FailedAt: &nodes[offset]}
		}

		if nodes[offset].GetType() == "PUNCTUATOR" && nodes[offset].GetImage() == "," {
			offset++
		}
	}

	return types.Result{Node: &node, Start: start, End: offset + 1}
}
