package rules

import (
	"go/ast"

	types "github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

type DotNotation struct {
	types.BaseNode
	Accessors []types.Node `json:"accessors"`
}

func (d *DotNotation) ToGoAst() ast.Node {
	// Start with the first identifier
	var expr ast.Expr
	expr = ast.NewIdent(d.Accessors[0].GetImage())

	// Chain the accesses
	for _, access := range d.Accessors[1:] {
		selector := &ast.SelectorExpr{
			X:   expr,
			Sel: ast.NewIdent(access.GetImage()),
		}

		expr = selector
	}

	return expr
}

func (d *DotNotation) StaticAnalysis(posh *types.PoshFile) {
	// if the first accessor is an identifier and not in the environment
	// AND is one of the built-in libraries, we need to add it to the environment
	if d.Accessors[0].GetType() == "IDENTIFIER" {
		image := d.Accessors[0].GetImage()
		if _, ok := posh.Environment.Get(image); !ok {
			if _, ok := StdImports[image]; ok {
				posh.StdImports[image] = true
			}
		}
	}
}

func MatchDotNotation(nodes []types.Node, offset int) types.Result {
	start := offset

	if nodes[offset].GetType() != "IDENTIFIER" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	node := DotNotation{
		BaseNode: types.BaseNode{
			Type: "DOT_NOTATION",
		},
		Accessors: []types.Node{nodes[offset]},
	}

	// if there are no dots, this is not a dot notation
	if offset+1 >= len(nodes) || nodes[offset+1].GetType() != "PUNCTUATOR" || nodes[offset+1].GetImage() != "." {
		return types.Result{FailedAt: &nodes[offset+1]}
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
