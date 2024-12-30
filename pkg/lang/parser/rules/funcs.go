package rules

import (
	"github.com/pouya-eghbali/alien-go/pkg/lang/parser/types"
)

type Param struct {
	types.BaseNode
	Identifier types.Node `json:"identifier"`
	ParamType  types.Node `json:"paramType"`
}

type Parameters struct {
	types.BaseNode
	Params []Param `json:"params"`
}

type FunctionBody struct {
	types.BaseNode
	Content []types.Node `json:"content"`
}

// TODO: Needs plug and unplug
type Function struct {
	types.BaseNode
	Identifier types.Node    `json:"identifier"`
	ReturnType types.Node    `json:"returnType"`
	Params     *Parameters   `json:"params"`
	Body       *FunctionBody `json:"body"`
}

func MatchFunctionParams(nodes []types.Node, offset int) types.Result {
	start := offset

	if nodes[offset].GetType() != "PUNCTUATOR" || nodes[offset].GetImage() != "(" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	offset++

	node := Parameters{
		BaseNode: types.BaseNode{
			Type: "PARAMETERS",
		},
	}

	for {
		if nodes[offset].GetType() == "PUNCTUATOR" && nodes[offset].GetImage() == ")" {
			break
		}

		if nodes[offset].GetType() != "IDENTIFIER" {
			return types.Result{FailedAt: &nodes[offset]}
		}

		param := Param{
			BaseNode: types.BaseNode{
				Type: "PARAM",
			},
			Identifier: nodes[offset],
		}

		offset++

		if nodes[offset].GetType() != "IDENTIFIER" {
			return types.Result{FailedAt: &nodes[offset]}
		}

		param.ParamType = nodes[offset]
		node.Params = append(node.Params, param)
		offset++

		if nodes[offset].GetType() == "PUNCTUATOR" && nodes[offset].GetImage() == "," {
			offset++
		}
	}

	return types.Result{Node: &node, Start: start, End: offset + 1}
}

func MatchFunctionBody(nodes []types.Node, offset int) types.Result {
	start := offset

	if nodes[offset].GetType() != "PUNCTUATOR" || nodes[offset].GetImage() != "{" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	offset++

	node := FunctionBody{
		BaseNode: types.BaseNode{
			Type: "FUNCTION_BODY",
		},
	}

	for {
		if nodes[offset].GetType() == "PUNCTUATOR" && nodes[offset].GetImage() == "}" {
			break
		}

		if res := MatchAssignment(nodes, offset); res.End > res.Start {
			node.Content = append(node.Content, res.Node)
			offset = res.End
		} else if res := MatchFunctionCall(nodes, offset); res.End > res.Start {
			node.Content = append(node.Content, res.Node)
			offset = res.End
		} else {
			return types.Result{FailedAt: &nodes[offset]}
		}
	}

	return types.Result{Node: &node, Start: start, End: offset + 1}
}

func MatchFunction(nodes []types.Node, offset int) types.Result {
	start := offset

	if nodes[offset].GetType() != "KEYWORD" || nodes[offset].GetImage() != "fn" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	offset++

	if nodes[offset].GetType() != "IDENTIFIER" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	node := Function{
		BaseNode: types.BaseNode{
			Type: "FUNCTION",
		},
		Identifier: nodes[offset],
	}

	offset++

	res := MatchFunctionParams(nodes, offset)
	if res.End == res.Start {
		return types.Result{FailedAt: res.FailedAt}
	}

	node.Params = res.Node.(*Parameters)
	offset = res.End

	if node.Identifier.GetImage() != "main" {
		if nodes[offset].GetType() != "IDENTIFIER" {
			return types.Result{FailedAt: &nodes[offset]}
		}

		node.ReturnType = nodes[offset]
		offset++
	}

	res = MatchFunctionBody(nodes, offset)
	if res.End == res.Start {
		return types.Result{FailedAt: res.FailedAt}
	}

	node.Body = res.Node.(*FunctionBody)
	return types.Result{Node: &node, Start: start, End: res.End}
}
