package types

type BaseNode struct {
	Type string `json:"type"`
}

func (n *BaseNode) Plug(env *Environment) {
	// do nothing
}

func (n *BaseNode) UnPlug(env *Environment) {
	// do nothing
}

func (n *BaseNode) GetPos() *Pos {
	return nil
}

func (n *BaseNode) GetImage() string {
	return ""
}

func (n *BaseNode) GetType() string {
	return n.Type
}
