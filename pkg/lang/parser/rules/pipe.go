package rules

import (
	"go/ast"
	"go/token"

	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

type Pipe struct {
	types.BaseNode
	Value FunctionCall `json:"value"`
}

func (n *Pipe) ToGoAst() ast.Node {
	return n.Value.ToGoAst()
}

func (n *Pipe) CollectTopLevelAssignments(posh *types.PoshFile) {
	n.Value.CollectTopLevelAssignments(posh)
}

type RunContext struct {
	types.BaseNode
}

func (n *RunContext) ToGoAst() ast.Node {
	// this is &exec.RunContext{}
	return &ast.UnaryExpr{
		Op: token.AND,
		X: &ast.CompositeLit{
			Type: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "exec"},
				Sel: &ast.Ident{Name: "RunContext"},
			},
		},
	}
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
			call := res.Node.(*FunctionCall)
			// prepend the args to the call args
			call.Args = append(args, call.Args...)
			args = []types.Node{call}
			node.Value = *call
			offset = res.End
		} else {
			break
		}
	}

	return types.Result{Node: &node, Start: start, End: offset}
}
