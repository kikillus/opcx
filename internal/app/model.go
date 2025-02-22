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
	activeNode          opc.NodeDef
	state               ui.ViewState
	width               int
	height              int

	connectionView ConnectionViewModel
	browseView     BrowseViewModel
	detailsView    DetailsViewModel
	recursiveView RecursiveViewModel
}

type ConnectionViewModel struct {

	connectionTextInput textinput.Model
	err error
}

type BrowseViewModel struct {
	viewport *viewport.Model
	activeNode opc.NodeDef
	nav *ui.Navigation
	err error
}

type DetailsViewModel struct {
	viewport *viewport.Model
	activeNode opc.NodeDef
	err error
}

type RecursiveViewModel struct {
	viewport *viewport.Model
	activeNode opc.NodeDef
	nav *ui.Navigation
	err error
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
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}
