package rules

import (
	"go/ast"
	"go/token"

	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

type Param struct {
	types.BaseNode
	Identifier types.Node `json:"identifier"`
	ParamType  types.Node `json:"paramType"`
}

type Parameters struct {
	types.BaseNode
	Params []Param `json:"params"`
}

type FunctionBody struct {
	types.BaseNode
	Content []types.Node `json:"content"`
}

func (n *FunctionBody) ToGoAst() ast.Node {
	body := []ast.Stmt{}
	for _, content := range n.Content {
		body = append(body, content.ToGoStatementAst())
	}

	return &ast.BlockStmt{
		List: body,
	}
}

func (n *FunctionBody) CollectTopLevelAssignments(posh *types.PoshFile) {
	for _, content := range n.Content {
		content.CollectTopLevelAssignments(posh)
	}
}

// TODO: Needs plug and unplug
type Function struct {
	types.BaseNode
	Identifier types.Node    `json:"identifier"`
	ReturnType *types.Node   `json:"returnType"`
	Params     *Parameters   `json:"params"`
	Body       *FunctionBody `json:"body"`
}

func getFlagVarName(paramType string) string {
	switch paramType {
	case "string":
		return "StringVar"
	case "int":
		return "IntVar"
	case "bool":
		return "BoolVar"
	default:
		return "StringVar"
	}
}

func getFlagDefaultValue(paramType string) ast.Expr {
	switch paramType {
	case "string":
		return &ast.BasicLit{
			Kind:  token.STRING,
			Value: "\"\"",
		}
	case "int":
		return &ast.BasicLit{
			Kind:  token.INT,
			Value: "0",
		}
	case "bool":
		return &ast.BasicLit{
			Kind:  token.STRING,
			Value: "false",
		}
	default:
		return &ast.BasicLit{
			Kind:  token.STRING,
			Value: "\"\"",
		}
	}
}

func (n *Function) ToGoAst() ast.Node {
	funcType := &ast.FuncType{}

	if n.ReturnType != nil && (*n.ReturnType).GetImage() != "void" {
		funcType.Results = &ast.FieldList{
			List: []*ast.Field{
				{
					Type: (*n.ReturnType).ToGoAst().(ast.Expr),
				},
			},
		}
	}

	body := n.Body.ToGoAst().(*ast.BlockStmt)

	if n.Identifier.GetImage() == "main" {
		// main function params should be turned into command line arguments
		// using flag package and added to the main function body:
		// fn main(name string, age int) should be turned into:
		// func main() {
		//     var name string
		//     flag.StringVar(&name, "name", "", "")
		//     var age int
		//     age := flag.IntVar(&age, "age", "", "")
		//     flag.Parse()
		// }
		// parameter types should be used to determine the type of the flag

		body.List = append([]ast.Stmt{&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "flag"},
					Sel: &ast.Ident{Name: "Parse"},
				},
			},
		}}, body.List...)

		for i := len(n.Params.Params) - 1; i >= 0; i-- {
			param := n.Params.Params[i]
			body.List = append([]ast.Stmt{
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{param.Identifier.ToGoAst().(*ast.Ident)},
								Type:  param.ParamType.ToGoAst().(ast.Expr),
							},
						},
					},
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "flag"},
							Sel: &ast.Ident{Name: getFlagVarName(param.ParamType.GetImage())},
						},
						Args: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X:  &ast.Ident{Name: param.Identifier.GetImage()},
							},
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"" + param.Identifier.GetImage() + "\"",
							},
							getFlagDefaultValue(param.ParamType.GetImage()),
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"\"",
							},
						},
					},
				},
			}, body.List...)
		}

	} else {
		params := []*ast.Field{}

		for _, param := range n.Params.Params {
			params = append(params, &ast.Field{
				Names: []*ast.Ident{param.Identifier.ToGoAst().(*ast.Ident)},
				Type:  param.ParamType.ToGoAst().(ast.Expr),
			})
		}

		funcType.Params = &ast.FieldList{
			List: params,
		}
	}

	return &ast.FuncDecl{
		Name: n.Identifier.ToGoAst().(*ast.Ident),
		Type: funcType,
		Body: body,
	}
}

func (n *Function) CollectTopLevelAssignments(posh *types.PoshFile) {
	posh.Environment.PushScope()

	for _, param := range n.Params.Params {
		posh.Environment.Set(param.Identifier.GetImage(), param.ParamType.GetImage())
	}

	n.Body.CollectTopLevelAssignments(posh)
	posh.Environment.PopScope()
}

func MatchFunctionParams(nodes []types.Node, offset int) types.Result {
	start := offset

	if nodes[offset].GetType() != "PUNCTUATOR" || nodes[offset].GetImage() != "(" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	offset++

	node := Parameters{
		BaseNode: types.BaseNode{
			Type: "PARAMETERS",
		},
	}

	for {
		if nodes[offset].GetType() == "PUNCTUATOR" && nodes[offset].GetImage() == ")" {
			break
		}

		if nodes[offset].GetType() != "IDENTIFIER" {
			return types.Result{FailedAt: &nodes[offset]}
		}

		param := Param{
			BaseNode: types.BaseNode{
				Type: "PARAM",
			},
			Identifier: nodes[offset],
		}

		offset++

		if nodes[offset].GetType() != "IDENTIFIER" {
			return types.Result{FailedAt: &nodes[offset]}
		}

		param.ParamType = nodes[offset]
		node.Params = append(node.Params, param)
		offset++

		if nodes[offset].GetType() == "PUNCTUATOR" && nodes[offset].GetImage() == "," {
			offset++
		}
	}

	return types.Result{Node: &node, Start: start, End: offset + 1}
}

func MatchFunctionBody(nodes []types.Node, offset int) types.Result {
	start := offset

	if nodes[offset].GetType() != "PUNCTUATOR" || nodes[offset].GetImage() != "{" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	offset++

	node := FunctionBody{
		BaseNode: types.BaseNode{
			Type: "FUNCTION_BODY",
		},
	}

	for {
		if nodes[offset].GetType() == "PUNCTUATOR" && nodes[offset].GetImage() == "}" {
			break
		}

		if res := MatchAssignment(nodes, offset); res.End > res.Start {
			node.Content = append(node.Content, res.Node)
			offset = res.End
		} else if res := MatchFunctionCall(nodes, offset); res.End > res.Start {
			node.Content = append(node.Content, res.Node)
			offset = res.End
		} else {
			return types.Result{FailedAt: &nodes[offset]}
		}
	}

	return types.Result{Node: &node, Start: start, End: offset + 1}
}

func MatchFunction(nodes []types.Node, offset int) types.Result {
	start := offset

	if nodes[offset].GetType() != "KEYWORD" || nodes[offset].GetImage() != "fn" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	offset++

	if nodes[offset].GetType() != "IDENTIFIER" {
		return types.Result{FailedAt: &nodes[offset]}
	}

	node := Function{
		BaseNode: types.BaseNode{
			Type: "FUNCTION",
		},
		Identifier: nodes[offset],
	}

	offset++

	res := MatchFunctionParams(nodes, offset)
	if res.End == res.Start {
		return types.Result{FailedAt: res.FailedAt}
	}

	node.Params = res.Node.(*Parameters)
	offset = res.End

	if node.Identifier.GetImage() != "main" {
		if nodes[offset].GetType() != "IDENTIFIER" {
			return types.Result{FailedAt: &nodes[offset]}
		}

		node.ReturnType = &nodes[offset]
		offset++
	}

	res = MatchFunctionBody(nodes, offset)
	if res.End == res.Start {
		return types.Result{FailedAt: res.FailedAt}
	}

	node.Body = res.Node.(*FunctionBody)
	return types.Result{Node: &node, Start: start, End: res.End}
}
