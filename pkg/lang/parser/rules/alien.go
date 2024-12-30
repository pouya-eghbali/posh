package rules

import (
	types "github.com/pouya-eghbali/alien-go/pkg/lang/parser/types"
)

type Alien struct {
	types.BaseNode
	Content []types.Node `json:"content"`
}

func MatchAlien(nodes []types.Node, offset int) types.Result {
	// Loop match the top-level nodes until no more matches are found
	start := offset

	node := Alien{
		BaseNode: types.BaseNode{
			Type: "ALIEN",
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
