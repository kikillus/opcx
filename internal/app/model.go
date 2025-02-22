package app

import (
	"opcx/internal/opc"
	"opcx/internal/ui"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	client              *opc.Client
	state               ui.ViewState
	width               int
	height              int

	connectionView ConnectionViewModel
	browseView     BrowseViewModel
	detailsView    DetailsViewModel
	recursiveView RecursiveViewModel
	monitorView   MonitorViewModel

	monitoredNodes []opc.NodeDef
}

type (
	errMsg error
)

func InitialModel() model {
	vp := viewport.New(0, 0)
	vp.YPosition = 3
	return model{
		state: ui.ViewStateConnection,
		connectionView: NewConnectionViewModel(),
		browseView: NewBrowseViewModel(),
		detailsView: NewDetailsViewModel(),
		recursiveView: NewRecursiveViewModel(),
		monitorView: NewMonitorViewModel(),
		monitoredNodes: make([]opc.NodeDef, 0),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}
