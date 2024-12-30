package rules

import (
	"go/ast"
	"go/token"

	types "github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

type Posh struct {
	types.BaseNode
	Content []types.Node `json:"content"`
}

var stdImports = []string{
	"github.com/pouya-eghbali/posh/pkg/exec",
	"github.com/pouya-eghbali/posh/pkg/io",
	"flag",
}

func (n *Posh) ToGoAst() ast.Node {
	posh := types.NewPoshFile()
	decls := []ast.Decl{}
	impSpecs := []ast.Spec{}

	// find all imports first
	for _, node := range n.Content {
		if node.GetType() == "IMPORT" {
			impDecl := node.ToGoAst().(ast.Decl)
			impSpecs = append(impSpecs, impDecl.(*ast.GenDecl).Specs...)
		}
	}

	// Add standard imports
	for _, imp := range stdImports {
		impSpecs = append(impSpecs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "\"" + imp + "\"",
			},
		})
	}

	// Add the import declaration
	decls = append(decls, &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: impSpecs,
	})

	// collect top-level assignments
	for _, node := range n.Content {
		node.CollectTopLevelAssignments(posh)
	}

	// add the value specs to the decls
	if len(posh.TopLevelAssignments) > 0 {
		decls = append(decls, &ast.GenDecl{
			Tok:   token.VAR,
			Specs: posh.TopLevelAssignments,
		})
	}

	// add the function declarations
	for _, node := range n.Content {
		if node.GetType() == "FUNCTION" {
			decls = append(decls, node.ToGoAst().(ast.Decl))
		}
	}

	return &ast.File{
		Name:  &ast.Ident{Name: "main"},
		Decls: decls,
	}
}

func MatchPosh(nodes []types.Node, offset int) types.Result {
	// Loop match the top-level nodes until no more matches are found
	start := offset

	node := Posh{
		BaseNode: types.BaseNode{
			Type: "POSH",
		},
	}

	for {
		if offset >= len(nodes) {
			break
		}

		var funRes types.Result
		var impRes types.Result

		// Match function
		if funRes = MatchFunction(nodes, offset); funRes.End > funRes.Start {
			node.Content = append(node.Content, funRes.Node)
			offset = funRes.End
			continue
		}

		// Match import
		if impRes = MatchImport(nodes, offset); impRes.End > impRes.Start {
			node.Content = append(node.Content, impRes.Node)
			offset = impRes.End
			continue
		}

		funPos := (*funRes.FailedAt).GetPos()
		impPos := (*impRes.FailedAt).GetPos()

		// Match error
		// return whichever has a bigger offset
		if funPos.Line > impPos.Line || (funPos.Line == impPos.Line && funPos.Column > impPos.Column) {
			return types.Result{FailedAt: funRes.FailedAt}
		} else {
			return types.Result{FailedAt: impRes.FailedAt}
		}
	}

	return types.Result{Node: &node, Start: start, End: offset}
}
