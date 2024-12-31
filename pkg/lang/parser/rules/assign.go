package rules

import (
	"go/ast"
	"go/token"

	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

// TODO: Needs plug and unplug
type Assignment struct {
	types.BaseNode
	Identifier types.Node `json:"identifier"`
	Value      types.Node `json:"value"`
	ValueType  string     `json:"valueType"`
}

func (a *Assignment) ToGoAst() ast.Node {
	var value ast.Expr

	if a.Value != nil && a.Value.GetType() == "PIPE" {
		// we need to add .Wait().ToString() to the end of the pipe
		value = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   a.Value.ToGoAst().(ast.Expr),
				Sel: &ast.Ident{Name: "Wait"},
			},
			Args: []ast.Expr{},
		}

		value = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   value,
				Sel: &ast.Ident{Name: "ToString"},
			},
			Args: []ast.Expr{},
		}
	} else if a.Value != nil {
		value = a.Value.ToGoAst().(ast.Expr)
	}

	return &ast.AssignStmt{
		Lhs: []ast.Expr{a.Identifier.ToGoAst().(ast.Expr)},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{value},
	}
}

func (a *Assignment) ToGoStatementAst() ast.Stmt {
	return a.ToGoAst().(*ast.AssignStmt)
}

func (a *Assignment) StaticAnalysis(posh *types.PoshFile) {
	posh.Environment.Set(a.Identifier.GetImage(), "unknown")
	a.Value.StaticAnalysis(posh)
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
