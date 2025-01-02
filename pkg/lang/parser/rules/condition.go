package rules

import (
	"go/ast"

	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

type Elif struct {
	types.BaseNode
	Condition types.Node `json:"condition"`
	Body      types.Node `json:"body"`
}

type Else struct {
	types.BaseNode
	Body types.Node `json:"body"`
}

type IfStatement struct {
	types.BaseNode
	Condition types.Node `json:"condition"`
	Body      types.Node `json:"body"`
	Elifs     []Elif     `json:"else_if"`
	Else      Else       `json:"else"`
}

func (n *IfStatement) ToGoStatementAst() ast.Stmt {
	ifNode := ast.IfStmt{
		Cond: n.Condition.ToGoAst().(ast.Expr),
		Body: n.Body.ToGoAst().(*ast.BlockStmt),
	}

	currentIfNode := &ifNode

	for _, elif := range n.Elifs {
		elifNode := &ast.IfStmt{
			Cond: elif.Condition.ToGoAst().(ast.Expr),
			Body: elif.Body.ToGoAst().(*ast.BlockStmt),
		}
		currentIfNode.Else = elifNode
		currentIfNode = elifNode
	}

	if n.Else.Body != nil {
		currentIfNode.Else = n.Else.Body.ToGoAst().(*ast.BlockStmt)
	}

	return &ifNode
}

func MatchIfStatement(nodes []types.Node, offset int) types.Result {
	start := offset

	// We are looking for the following:
	// IF BOOLEAN BODY (ELSEIF BOOLEAN|LOGICAL BODY)* (ELSE BODY)?

	// try to match IF
	if nodes[offset].GetType() != "KEYWORD" || nodes[offset].GetImage() != "if" {
		return types.Result{FailedAt: &nodes[offset]}
	}
	offset++

	// create the node
	node := IfStatement{
		BaseNode: types.BaseNode{
			Type: "IF",
		},
	}

	// try to match BOOLEAN or LOGICAL
	if res := MatchLogical(nodes, offset); res.End > res.Start {
		offset = res.End
		node.Condition = res.Node
	} else if res := MatchBoolean(nodes, offset); res.End > res.Start {
		offset = res.End
		node.Condition = res.Node
	} else {
		return types.Result{FailedAt: &nodes[offset]}
	}

	// try to match BODY; we can reuse MatchFunctionBody here
	if res := MatchFunctionBody(nodes, offset); res.End <= res.Start {
		return types.Result{FailedAt: &nodes[offset]}
	} else {
		offset = res.End
		node.Body = res.Node
	}

	// try to match (ELSEIF BOOLEAN|LOGICAL BODY)*
	// PoSH uses "elif" instead of "elseif" or "else if"
	for {
		if nodes[offset].GetType() != "KEYWORD" || nodes[offset].GetImage() != "elif" {
			break
		}
		offset++

		elifNode := Elif{
			BaseNode: types.BaseNode{
				Type: "ELIF",
			},
		}

		// try to match BOOLEAN or LOGICAL
		if res := MatchLogical(nodes, offset); res.End > res.Start {
			offset = res.End
			elifNode.Condition = res.Node
		} else if res := MatchBoolean(nodes, offset); res.End > res.Start {
			offset = res.End
			elifNode.Condition = res.Node
		} else {
			return types.Result{FailedAt: &nodes[offset]}
		}

		// try to match BODY
		if res := MatchFunctionBody(nodes, offset); res.End <= res.Start {
			return types.Result{FailedAt: &nodes[offset]}
		} else {
			offset = res.End
			elifNode.Body = res.Node
		}

		node.Elifs = append(node.Elifs, elifNode)
	}

	// try to match (ELSE BODY)?
	if nodes[offset].GetType() == "KEYWORD" && nodes[offset].GetImage() == "else" {
		elseNode := Else{
			BaseNode: types.BaseNode{
				Type: "ELSE",
			},
		}
		offset++

		// try to match BODY
		if res := MatchFunctionBody(nodes, offset); res.End <= res.Start {
			return types.Result{FailedAt: &nodes[offset]}
		} else {
			offset = res.End
			elseNode.Body = res.Node
			node.Else = elseNode
		}
	}

	return types.Result{Node: &node, Start: start, End: offset}
}
