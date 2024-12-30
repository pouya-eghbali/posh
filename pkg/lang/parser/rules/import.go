package rules

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/pouya-eghbali/alien-go/pkg/lang/parser/types"
)

type ImportItem struct {
	types.BaseNode
	Name  types.Node `json:"name"`
	Alias types.Node `json:"alias"`
}

// TODO: Needs plug and unplug
type Import struct {
	types.BaseNode
	Module  types.Node    `json:"module"`
	Imports []*ImportItem `json:"imports"`
}

func (n *Import) ToGoAst() ast.Node {
	// create a list of imports
	if len(n.Imports) == 1 && n.Imports[0].Name.GetImage() == "*" {
		// import all
		return &ast.GenDecl{
			Tok: token.IMPORT,
			Specs: []ast.Spec{
				&ast.ImportSpec{
					Name: n.Imports[0].Alias.ToGoAst().(*ast.Ident),
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: n.Module.GetImage(),
					},
				},
			},
		}
	}

	specs := []ast.Spec{}
	for _, imp := range n.Imports {
		spec := &ast.ImportSpec{
			Name: imp.Alias.ToGoAst().(*ast.Ident),
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: n.Module.GetImage(),
			},
		}
		specs = append(specs, spec)
	}

	return &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: specs,
	}
}

func (n *Import) CollectTopLevelAssignments(alien *types.AlienFile) {
	// for each import item as item, create an assignment
	// e.g. from "fmt" import Println as fmtPrintln then
	// fmtPrintln = fmt.Println
	node := ast.ValueSpec{}

	for _, imp := range n.Imports {
		if imp.Name.GetImage() != "*" && imp.Alias == nil {
			node.Names = append(node.Names, imp.Name.ToGoAst().(*ast.Ident))
			node.Values = append(node.Values, &ast.SelectorExpr{
				X:   &ast.Ident{Name: importPathToImportName(n.Module.GetImage())},
				Sel: imp.Name.ToGoAst().(*ast.Ident),
			})
		}
	}

	if len(node.Names) > 0 {
		alien.TopLevelAssignments = append(alien.TopLevelAssignments, &node)
	}
}

func importPathToImportName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

func MatchImportAllAsItem(nodes []types.Node, offset int) types.Result {
	// We are looking for the following:
	// * as <identifier>
	start := offset

	// try to match "*"
	if nodes[offset].GetType() != "PUNCTUATOR" || nodes[offset].GetImage() != "*" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	node := ImportItem{
		BaseNode: types.BaseNode{
			Type: "IMPORT_ITEM",
		},
		Name: nodes[offset],
	}

	offset++

	// try to match "as"
	if nodes[offset].GetType() != "KEYWORD" || nodes[offset].GetImage() != "as" {
		return types.Result{FailedAt: &nodes[offset]}
	}
	offset++

	// try to match <identifier>
	if nodes[offset].GetType() != "IDENTIFIER" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	node.Alias = nodes[offset]
	offset++

	return types.Result{Node: &node, Start: start, End: offset}

}

func MatchImportItem(nodes []types.Node, offset int) types.Result {
	// We are looking for the following:
	// identifier (as identifier)?)
	start := offset

	// try to match <identifier>
	if nodes[offset].GetType() != "IDENTIFIER" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	node := ImportItem{
		BaseNode: types.BaseNode{
			Type: "IMPORT_ITEM",
		},
		Name: nodes[offset],
	}

	offset++

	// try to match "as"
	if nodes[offset].GetType() == "KEYWORD" && nodes[offset].GetImage() == "as" {
		offset++

		// try to match <identifier>
		if nodes[offset].GetType() != "IDENTIFIER" {
			return types.Result{FailedAt: &nodes[offset]}
		}

		node.Alias = nodes[offset]
		offset++
	}

	return types.Result{Node: &node, Start: start, End: offset}
}

func MatchImport(nodes []types.Node, offset int) types.Result {
	start := offset

	// We are looking for one of the following:
	// from <string> import (identifier (as identifier)?) (, identifier (as identifier)?)*
	// from <string> import * as <identifier>

	imp := Import{
		BaseNode: types.BaseNode{
			Type: "IMPORT",
		},
	}

	// try to match "from"
	if nodes[offset].GetType() != "KEYWORD" || nodes[offset].GetImage() != "from" {
		return types.Result{FailedAt: &nodes[offset]}
	}
	offset++

	// try to match <string>
	if nodes[offset].GetType() != "STRING" {
		return types.Result{FailedAt: &nodes[offset]}
	}
	imp.Module = nodes[offset]
	offset++

	// try to match "import"
	if nodes[offset].GetType() != "KEYWORD" || nodes[offset].GetImage() != "import" {
		return types.Result{FailedAt: &nodes[offset]}
	}
	offset++

	// try to match importAllAsItem
	if res := MatchImportAllAsItem(nodes, offset); res.End > res.Start {
		imp.Imports = append(imp.Imports, res.Node.(*ImportItem))
		offset = res.End
	} else {
		for {
			if res := MatchImportItem(nodes, offset); res.End > res.Start {
				imp.Imports = append(imp.Imports, res.Node.(*ImportItem))
				offset = res.End
			} else {
				break
			}

			if nodes[offset].GetType() != "PUNCTUATOR" || nodes[offset].GetImage() != "," {
				break
			}
			offset++
		}
	}

	return types.Result{Node: &imp, Start: start, End: offset}
}
