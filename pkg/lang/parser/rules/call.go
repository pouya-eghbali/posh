package rules

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

type Flag struct {
	types.BaseNode
	Identifier types.Node `json:"identifier"`
	DashCount  int        `json:"dashCount"`
}

func (n *Flag) ToGoAst() ast.Node {
	// treat the flag as a string
	flagStr := strings.Repeat("-", n.DashCount) + n.Identifier.GetImage()
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: fmt.Sprintf(`"%s"`, flagStr),
	}
}

type FunctionCall struct {
	types.BaseNode
	Callable types.Node   `json:"callable"`
	Args     []types.Node `json:"args"`
}

func (n *FunctionCall) ToGoAst() ast.Node {
	args := []ast.Expr{}
	for _, arg := range n.Args {
		args = append(args, arg.ToGoAst().(ast.Expr))
	}

	return &ast.CallExpr{
		Fun:  n.Callable.ToGoAst().(ast.Expr),
		Args: args,
	}
}

func (n *FunctionCall) StaticAnalysis(posh *types.PoshFile) {
	n.Callable.StaticAnalysis(posh)

	if n.Callable.GetType() == "IDENTIFIER" {
		image := n.Callable.GetImage()
		if _, ok := posh.Environment.Get(image); !ok {
			// We need to add {identifier} := exec.ExternalCommand("{identifier}")
			posh.TopLevelAssignments = append(posh.TopLevelAssignments, &ast.ValueSpec{
				Names: []*ast.Ident{{Name: image}},
				Values: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "exec"},
							Sel: &ast.Ident{Name: "ExternalCommand"},
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: fmt.Sprintf(`"%s"`, image),
							},
						},
					},
				},
			})
		}
	}

	for _, arg := range n.Args {
		arg.StaticAnalysis(posh)
	}
}

func (n *FunctionCall) ToGoStatementAst() ast.Stmt {
	return &ast.ExprStmt{
		X: n.ToGoAst().(ast.Expr),
	}
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
