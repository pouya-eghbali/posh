package types

type Pos struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Node interface {
	GetPos() *Pos
	GetImage() string
	GetType() string
	Plug(*Environment)
	UnPlug(*Environment)
}

type Result struct {
	Node     Node
	FailedAt *Node
	Start    int
	End      int
}
