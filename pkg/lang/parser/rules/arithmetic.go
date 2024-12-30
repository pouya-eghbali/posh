package rules

import (
	"go/ast"
	"go/token"

	"github.com/pouya-eghbali/alien-go/pkg/lang/parser/types"
)

type ArithmeticNode struct {
	types.BaseNode
	Lhs types.Node `json:"lhs"`
	Op  types.Node `json:"op"`
	Rhs types.Node `json:"rhs"`
}

var tokenMap = map[string]token.Token{
	"+": token.ADD,
	"-": token.SUB,
	"*": token.MUL,
	"/": token.QUO,
}

func (a *ArithmeticNode) ToGoAst() ast.Node {
	var lhs, rhs ast.Expr

	if a.Lhs != nil {
		lhs = a.Lhs.ToGoAst().(ast.Expr)
	}

	if a.Rhs != nil {
		rhs = a.Rhs.ToGoAst().(ast.Expr)
	}

	return &ast.BinaryExpr{
		X:  lhs,
		Op: tokenMap[a.Op.GetImage()],
		Y:  rhs,
	}
}

func isOperator(value string) bool {
	return value == "+" || value == "-" || value == "*" || value == "/"
}

func MatchWrappedArithmetic(nodes []types.Node, offset int) types.Result {
	start := offset

	if nodes[offset].GetType() != "PUNCTUATOR" || nodes[offset].GetImage() != "(" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	offset++

	var res types.Result
	if res = MatchArithmetic(nodes, offset); res.End > res.Start {
		node := ArithmeticNode{
			BaseNode: types.BaseNode{
				Type: "WRAPPED_ARITHMETIC",
			},
			Lhs: res.Node.(*ArithmeticNode).Lhs,
			Op:  res.Node.(*ArithmeticNode).Op,
			Rhs: res.Node.(*ArithmeticNode).Rhs,
		}

		if nodes[res.End].GetType() != "PUNCTUATOR" || nodes[res.End].GetImage() != ")" {
			return types.Result{FailedAt: &nodes[res.End]}
		}

		return types.Result{Node: &node, Start: start, End: res.End + 1}
	}

	return types.Result{FailedAt: res.FailedAt}
}

func MatchArithmetic(nodes []types.Node, offset int) types.Result {
	start := offset

	// We should recursively match the following:
	// NUMERIC [OPERATOR NUMERIC]+
	// We don't care about the order of the operators
	// since Go will handle the precedence for us

	// try to match NUMERIC
	var res types.Result
	if res = MatchNumeric(nodes, offset); res.End > res.Start {
		offset = res.End
	} else {
		return types.Result{FailedAt: res.FailedAt}
	}

	node := ArithmeticNode{
		BaseNode: types.BaseNode{
			Type: "ARITHMETIC",
		},
		Lhs: res.Node,
	}

	// if there are no operators, this is not an arithmetic expression
	if offset+1 >= len(nodes) || nodes[offset+1].GetType() != "PUNCTUATOR" || !isOperator(nodes[offset+1].GetImage()) {
		return types.Result{FailedAt: &nodes[offset+1]}
	}

	for {
		if nodes[offset].GetType() != "PUNCTUATOR" || !isOperator(nodes[offset].GetImage()) {
			break
		}

		node.Op = nodes[offset]
		offset++

		if res = MatchNumeric(nodes, offset); res.End > res.Start {
			node.Rhs = res.Node
			offset = res.End
		} else {
			return types.Result{FailedAt: res.FailedAt}
		}

		// if there's more operators, we need to create a new node
		// and set the current node as the LHS
		if offset+1 < len(nodes) && nodes[offset+1].GetType() == "PUNCTUATOR" && isOperator(nodes[offset+1].GetImage()) {
			node = ArithmeticNode{
				BaseNode: types.BaseNode{
					Type: "ARITHMETIC",
				},
				Lhs: &node,
			}
		}
	}

	return types.Result{Node: &node, Start: start, End: offset}
}
