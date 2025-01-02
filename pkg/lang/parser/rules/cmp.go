package rules

import (
	"go/ast"
	"go/token"

	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

type ComparisonNode struct {
	types.BaseNode
	Lhs types.Node `json:"lhs"`
	Op  types.Node `json:"op"`
	Rhs types.Node `json:"rhs"`
}

var cmpTokenMap = map[string]token.Token{
	">":  token.GTR,
	"<":  token.LSS,
	">=": token.GEQ,
	"<=": token.LEQ,
	"==": token.EQL,
	"!=": token.NEQ,
}

func (a *ComparisonNode) ToGoAst() ast.Node {
	var lhs, rhs ast.Expr

	if a.Lhs != nil {
		lhs = a.Lhs.ToGoAst().(ast.Expr)
	}

	if a.Rhs != nil {
		rhs = a.Rhs.ToGoAst().(ast.Expr)
	}

	op := cmpTokenMap[a.Op.GetImage()]

	return &ast.BinaryExpr{
		X:  lhs,
		Op: op,
		Y:  rhs,
	}
}

func isCmpOperator(value string) bool {
	return value == ">" || value == "<"
}

func MatchComparison(nodes []types.Node, offset int) types.Result {
	start := offset

	// We should recursively match the following:
	// NUMERIC (COMPARISON_OPERATOR NUMERIC)*
	// In PoSH, a < b < c is valid and is equivalent to (a < b) && (b < c)

	var res types.Result
	if res = MatchNumeric(nodes, offset); res.End <= res.Start {
		return types.Result{FailedAt: &nodes[offset]}
	}
	offset = res.End

	cmp := ComparisonNode{
		BaseNode: types.BaseNode{
			Type: "COMPARISON",
		},
		Lhs: res.Node,
	}

	cmpList := []ComparisonNode{}

	for offset < len(nodes) {
		if nodes[offset].GetType() != "PUNCTUATOR" {
			break
		}

		// We should check if the operator is a comparison operator
		// Some operators are two characters long, so we should check for that
		image := nodes[offset].GetImage()
		deltaOffset := 1
		if image == "!" {
			if nodes[offset+1].GetImage() != "=" {
				break
			}
			deltaOffset = 2
			image = "!="
		} else if image == "=" {
			if nodes[offset+1].GetImage() != "=" {
				break
			}
			deltaOffset = 2
			image = "=="
		} else if !isCmpOperator(image) {
			break
		}

		cmp.Op = &types.TokenNode{
			BaseNode: types.BaseNode{
				Type: "COMPARISON_OPERATOR",
			},
			Image: image,
			Pos:   nodes[offset].GetPos(),
		}

		offset += deltaOffset

		// Now we should match the right hand side of the comparison
		if res = MatchNumeric(nodes, offset); res.End <= res.Start {
			return types.Result{FailedAt: &nodes[offset]}
		}
		offset = res.End
		cmp.Rhs = res.Node

		cmpList = append(cmpList, cmp)
		cmp = ComparisonNode{
			BaseNode: types.BaseNode{
				Type: "COMPARISON",
			},
			Lhs: res.Node,
		}
	}

	if len(cmpList) == 0 {
		return types.Result{FailedAt: &nodes[offset]}
	} else if len(cmpList) == 1 {
		return types.Result{Node: &cmpList[0], Start: start, End: offset}
	}

	// a < b < c should be equivalent to (a < b) && (b < c)
	// We should build a tree of AND nodes

	logicNode := Logical{
		BaseNode: types.BaseNode{
			Type: "LOGICAL",
		},
		Op: &types.TokenNode{
			BaseNode: types.BaseNode{
				Type: "LOGICAL_OPERATOR",
			},
			Image: "and",
			Pos:   cmpList[0].GetPos(),
		},
		Lhs: &cmpList[0],
	}

	for i := 1; i < len(cmpList); i++ {
		logicNode.Rhs = &Logical{
			BaseNode: types.BaseNode{
				Type: "LOGICAL",
			},
			Op: &types.TokenNode{
				BaseNode: types.BaseNode{
					Type: "LOGICAL_OPERATOR",
				},
				Image: "and",
				Pos:   cmpList[i].GetPos(),
			},
			Lhs: &cmpList[i],
		}
	}

	return types.Result{Node: &logicNode, Start: start, End: offset}
}
