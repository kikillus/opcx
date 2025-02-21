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
	nav                 *ui.Navigation
	activeNode          opc.NodeDef
	state               ui.ViewState
	width               int
	height              int
	connectionTextInput textinput.Model
	err                 error
	viewport            *viewport.Model
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
	return model{nav: ui.NewNavigation(), viewport: &vp, connectionTextInput: connectionText, state: ui.ViewStateConnection}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}
