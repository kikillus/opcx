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
	connectionText := textinput.New()
	connectionText.Placeholder = "opc.tcp://127.0.0.1:4840"
	connectionText.SetValue("opc.tcp://127.0.0.1:4840")
	connectionText.CharLimit = 128
	connectionText.Width = 50
	connectionText.Focus()
	vp := viewport.New(0, 0)
	vp.YPosition = 3
	return model{
		connectionView: NewConnectionViewModel(),
		browseView: NewBrowseViewModel(),
		detailsView: NewDetailsViewModel(),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}
