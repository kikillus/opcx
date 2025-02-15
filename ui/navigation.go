package ui

import (
	opcutil "opc-tui/opc/util"

	"github.com/gopcua/opcua/ua"
)

type Navigation struct {
	Path         []opcutil.NodeDef
	CurrentNodes []opcutil.NodeDef
	Cursor       int
}

func NewNavigation() *Navigation {
	rootNode := opcutil.NodeDef{
		NodeID: ua.NewNumericNodeID(0, 84),
	}

	return &Navigation{Path: []opcutil.NodeDef{rootNode}}
}

func (n *Navigation) Forward(node opcutil.NodeDef) {
	n.Path = append(n.Path, node)
}

func (n *Navigation) Backward() bool {
	if len(n.Path) <= 1 {
		return false
	}
	n.Path = n.Path[:len(n.Path)-1]
	return true
}

func (n *Navigation) CurrentNode() *opcutil.NodeDef {
	if len(n.Path) == 0 {
		return nil
	}
	return &n.Path[len(n.Path)-1]
}
