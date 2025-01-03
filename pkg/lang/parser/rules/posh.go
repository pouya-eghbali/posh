package rules

import (
	"go/ast"
	"go/token"

	types "github.com/pouya-eghbali/posh/pkg/lang/parser/types"
	"github.com/pouya-eghbali/posh/pkg/lang/parser/utils"
)

type Posh struct {
	types.BaseNode
	Content []types.Node `json:"content"`
}

var StdImports = map[string]string{
	"exec": "github.com/pouya-eghbali/posh/pkg/exec",
	"io":   "github.com/pouya-eghbali/posh/pkg/io",
	"flag": "flag",
}

func (n *Posh) CompileToGo(posh *types.PoshFile) error {
	decls := []ast.Decl{}
	impSpecs := []ast.Spec{}

	// perform self-analysis
	n.StaticAnalysis(posh)

	// find all imports first
	for _, node := range n.Content {
		if node.GetType() == "IMPORT" {
			impDecl := node.ToGoAst().(ast.Decl)
			impSpecs = append(impSpecs, impDecl.(*ast.GenDecl).Specs...)
		}
	}

	// Add standard imports from posh.StdImports
	for imp := range posh.StdImports {
		impSpecs = append(impSpecs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "\"" + StdImports[imp] + "\"",
			},
		})
	}

	// Add the import declaration
	decls = append(decls, &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: impSpecs,
	})

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

	file := &ast.File{
		Name:  &ast.Ident{Name: posh.Package},
		Decls: decls,
	}

	return utils.WritePoshFile(file, posh)
}

func (n *Posh) StaticAnalysis(posh *types.PoshFile) {
	// find all functions and add them to the environment
	for _, node := range n.Content {
		if node.GetType() == "FUNCTION" {
			// TODO: Environment should be a map of string to types.Export
			// TODO: Rename types.Export to something more meaningful
			posh.Environment.Set(node.(*Function).Identifier.GetImage(), "unknown")
		}
	}

	for _, node := range n.Content {
		node.StaticAnalysis(posh)
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
