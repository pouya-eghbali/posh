package types

import (
	"go/ast"
	"go/token"
)

type TokenNode struct {
	BaseNode
	Image string `json:"image"`
	Pos   *Pos   `json:"pos"`
}

func (n *TokenNode) GetPos() *Pos {
	return n.Pos
}

func (n *TokenNode) GetImage() string {
	return n.Image
}

func (n *TokenNode) ToGoAst() ast.Node {
	if n.Image == "true" {
		return &ast.Ident{Name: "true"}
	} else if n.Image == "false" {
		return &ast.Ident{Name: "false"}
	} else if n.Type == "STRING" {
		return &ast.BasicLit{Kind: token.STRING, Value: n.Image}
	} else if n.Type == "INTEGER" {
		return &ast.BasicLit{Kind: token.INT, Value: n.Image}
	} else if n.Type == "FLOAT" {
		return &ast.BasicLit{Kind: token.FLOAT, Value: n.Image}
	} else if n.Type == "IDENTIFIER" {
		return &ast.Ident{Name: n.Image}
	} else {
		return nil
	}
}
