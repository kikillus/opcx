package ui

import (
	"opcx/internal/opc"

	"github.com/gopcua/opcua/ua"
)

type Navigation struct {
	Path         []opc.NodeDef
	CurrentNodes []opc.NodeDef
	Cursor       int
	ActiveNode opc.NodeDef
}

func NewNavigation() *Navigation {
	rootNode := opc.NodeDef{
		NodeID: ua.NewNumericNodeID(0, 84),
	}

	return &Navigation{Path: []opc.NodeDef{rootNode}}
}

func (n *Navigation) Forward(node opc.NodeDef) {
	n.Path = append(n.Path, node)
}

func (n *Navigation) Backward() bool {
	if len(n.Path) <= 1 {
		return false
	}
	n.Path = n.Path[:len(n.Path)-1]
	return true
}

func (n *Navigation) CurrentNode() *opc.NodeDef {
	if len(n.Path) == 0 {
		return nil
	}
	return &n.Path[len(n.Path)-1]
}
