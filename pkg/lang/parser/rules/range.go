package rules

import (
	"go/ast"

	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

type Range struct {
	types.BaseNode
	Start types.Node  `json:"start"`
	Step  types.Node  `json:"step"`
	End   *types.Node `json:"end"`
}

func (a *Range) StaticAnalysis(posh *types.PoshFile) {
	a.Start.StaticAnalysis(posh)
	a.Step.StaticAnalysis(posh)

	if a.End != nil {
		(*a.End).StaticAnalysis(posh)
	}

	posh.StdImports["std"] = true
}

func (a *Range) ToGoAst() ast.Node {
	// return std.LazyRange(start, step, end...)
	r := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "std"},
			Sel: &ast.Ident{Name: "LazyRange"},
		},
		Args: []ast.Expr{
			a.Start.ToGoAst().(ast.Expr),
			a.Step.ToGoAst().(ast.Expr),
		},
	}

	if a.End != nil {
		r.Args = append(r.Args, (*a.End).ToGoAst().(ast.Expr))
	}

	return r
}

func MatchRange(nodes []types.Node, offset int) types.Result {
	start := offset

	// We are looking for the following:
	// NUMERIC,(,NUMERIC)?..NUMERIC?

	node := Range{
		BaseNode: types.BaseNode{
			Type: "RANGE",
		},
	}

	// Look for the start of the range
	if res := MatchNumeric(nodes, offset); res.End > res.Start {
		node.Start = res.Node
		offset = res.End
	} else {
		return types.Result{FailedAt: res.FailedAt}
	}

	// Look for the step of the range
	if nodes[offset].GetType() == "PUNCTUATOR" && nodes[offset].GetImage() == "," {
		offset++
		if res := MatchNumeric(nodes, offset); res.End > res.Start {
			node.Step = res.Node
			offset = res.End
		} else {
			return types.Result{FailedAt: res.FailedAt}
		}
	}

	// Look for ..
	for i := 0; i < 2; i++ {
		if nodes[offset].GetType() != "PUNCTUATOR" || nodes[offset].GetImage() != "." {
			return types.Result{FailedAt: &nodes[offset]}
		}
		offset++
	}

	// Look for the end of the range
	if res := MatchNumeric(nodes, offset); res.End > res.Start {
		node.End = &res.Node
		offset = res.End
	}

	return types.Result{Node: &node, Start: start, End: offset}
}
