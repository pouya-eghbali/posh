package rules

import "github.com/pouya-eghbali/alien-go/pkg/lang/parser/types"

type Pipe struct {
	types.BaseNode
	Calls []*FunctionCall `json:"calls"`
}

type RunContext struct {
	types.BaseNode
}

func MatchPipe(nodes []types.Node, offset int) types.Result {
	start := offset
	// We're looking for the following:
	// SimpleExpression (| FunctionCall)+

	node := Pipe{
		BaseNode: types.BaseNode{
			Type: "PIPE",
		},
	}

	args := []types.Node{&RunContext{}}

	// look for the simple expression and add it to the args
	if res := MatchSimpleExpression(nodes, offset); res.End > res.Start {
		args = append(args, res.Node)
		offset = res.End
	} else {
		return types.Result{FailedAt: res.FailedAt}
	}

	// if there's no pipe this is not a match
	if offset >= len(nodes) || nodes[offset].GetType() != "PUNCTUATOR" || nodes[offset].GetImage() != "|" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	// recursively look for (| FunctionCall)
	for {
		if nodes[offset].GetType() != "PUNCTUATOR" || nodes[offset].GetImage() != "|" {
			break
		}

		offset++

		if res := MatchFunctionCall(nodes, offset); res.End > res.Start {
			node.Calls = append(node.Calls, res.Node.(*FunctionCall))
			offset = res.End
		} else {
			break
		}
	}

	return types.Result{Node: &node, Start: start, End: offset}
}
