package rules

import (
	"go/ast"
	"go/token"

	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

// break and continue
type ForControl struct {
	types.BaseNode
	Op string `json:"op"`
}

func (n *ForControl) ToGoAst() ast.Node {
	if n.Op == "break" {
		return &ast.BranchStmt{
			Tok: token.BREAK,
		}
	} else {
		return &ast.BranchStmt{
			Tok: token.CONTINUE,
		}
	}
}

type ForBody struct {
	types.BaseNode
	Content []types.Node `json:"content"`
}

func (n *ForBody) ToGoAst() ast.Node {
	body := []ast.Stmt{}
	for _, content := range n.Content {
		body = append(body, content.ToGoStatementAst())
	}

	return &ast.BlockStmt{
		List: body,
	}
}

func (n *ForBody) StaticAnalysis(posh *types.PoshFile) {
	for _, content := range n.Content {
		content.StaticAnalysis(posh)
	}
}

func isForControl(s string) bool {
	return s == "break" || s == "continue"
}

func MatchForBody(nodes []types.Node, offset int) types.Result {
	start := offset

	if nodes[offset].GetType() != "PUNCTUATOR" || nodes[offset].GetImage() != "{" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	offset++

	node := ForBody{
		BaseNode: types.BaseNode{
			Type: "FUNCTION_BODY",
		},
	}

	for {
		if nodes[offset].GetType() == "PUNCTUATOR" && nodes[offset].GetImage() == "}" {
			offset++
			break
		}

		if res := MatchAssignment(nodes, offset); res.End > res.Start {
			node.Content = append(node.Content, res.Node)
			offset = res.End
		} else if res := MatchFunctionCall(nodes, offset); res.End > res.Start {
			node.Content = append(node.Content, res.Node)
			offset = res.End
		} else if res := MatchReturnStatement(nodes, offset); res.End > res.Start {
			node.Content = append(node.Content, res.Node)
			offset = res.End
		} else if res := MatchIfStatement(nodes, offset); res.End > res.Start {
			node.Content = append(node.Content, res.Node)
			offset = res.End
		} else if res := MatchForLoop(nodes, offset); res.End > res.Start {
			node.Content = append(node.Content, res.Node)
			offset = res.End
		} else if nodes[offset].GetType() == "KEYWORD" && isForControl(nodes[offset].GetImage()) {
			node.Content = append(node.Content, &ForControl{
				BaseNode: types.BaseNode{
					Type: "FOR_CONTROL",
				},
				Op: nodes[offset].GetImage(),
			})
			offset++
		} else {
			return types.Result{FailedAt: &nodes[offset]}
		}
	}

	return types.Result{Node: &node, Start: start, End: offset}
}

type ForLoop struct {
	types.BaseNode
	Variables []types.Node `json:"variables"`
	Iterable  types.Node   `json:"iterable"`
	Body      *ForBody     `json:"body"`
}

func (n *ForLoop) StaticAnalysis(posh *types.PoshFile) {
	n.Iterable.StaticAnalysis(posh)
	n.Body.StaticAnalysis(posh)
}

func (n *ForLoop) ToGoStatementAst() ast.Stmt {
	// this is the "k" part of "for k, v := range iterable"
	keyVar := &ast.Ident{Name: n.Variables[0].GetImage()}

	// this is the "v" part of "for k, v := range iterable"
	var valueVar ast.Expr
	if len(n.Variables) > 1 {
		valueVar = &ast.Ident{Name: n.Variables[1].GetImage()}
	}

	iterableExpr := n.Iterable.ToGoAst().(ast.Expr)
	bodyStmt := n.Body.ToGoAst().(*ast.BlockStmt)

	// Create the range statement
	return &ast.RangeStmt{
		Key:   keyVar,
		Value: valueVar,
		Tok:   token.DEFINE,
		X:     iterableExpr,
		Body:  bodyStmt,
	}
}

func MatchForLoop(nodes []types.Node, offset int) types.Result {
	start := offset

	// We are looking for the following:
	// FOR IDENTIFIER (, IDENTIFIER)? IN EXPRESSION BODY

	// try to match IF
	if nodes[offset].GetType() != "KEYWORD" || nodes[offset].GetImage() != "for" {
		return types.Result{FailedAt: &nodes[offset]}
	}
	offset++

	// create the node
	node := ForLoop{
		BaseNode: types.BaseNode{
			Type: "FOR",
		},
	}

	// try to match IDENTIFIER
	if nodes[offset].GetType() != "IDENTIFIER" {
		return types.Result{FailedAt: &nodes[offset]}
	}
	node.Variables = append(node.Variables, nodes[offset])
	offset++

	// try to match (, IDENTIFIER)?
	if nodes[offset].GetType() == "PUNCTUATOR" && nodes[offset].GetImage() == "," {
		offset++

		if nodes[offset].GetType() != "IDENTIFIER" {
			return types.Result{FailedAt: &nodes[offset]}
		}

		node.Variables = append(node.Variables, nodes[offset])
		offset++
	}

	// try to match IN
	if nodes[offset].GetType() != "KEYWORD" || nodes[offset].GetImage() != "in" {
		return types.Result{FailedAt: &nodes[offset]}
	}
	offset++

	// try to match EXPRESSION
	if res := MatchSimpleExpression(nodes, offset); res.End <= res.Start {
		return types.Result{FailedAt: &nodes[offset]}
	} else {
		offset = res.End
		node.Iterable = res.Node
	}

	// try to match BODY
	if res := MatchForBody(nodes, offset); res.End <= res.Start {
		return types.Result{FailedAt: &nodes[offset]}
	} else {
		offset = res.End
		node.Body = res.Node.(*ForBody)
	}

	return types.Result{Node: &node, Start: start, End: offset}
}
