package types

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
