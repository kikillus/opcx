package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ConnectionViewModel struct {
	connectionTextInput textinput.Model
	err                 error
}

func NewConnectionViewModel() ConnectionViewModel {
	ti := textinput.New()
	ti.Placeholder = "opc.tcp://127.0.0.1:4840"
	ti.SetValue("opc.tcp://127.0.0.1:4840")
	ti.CharLimit = 128
	ti.Width = 50
	ti.Focus()

	return ConnectionViewModel{
		connectionTextInput: ti,
	}
}


type TransitionConnectToBrowseMsg struct{
	endpoint string
}

func (m ConnectionViewModel) Update(msg tea.Msg) (ConnectionViewModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			newModel := m

			return newModel, func() tea.Msg {
					return TransitionConnectToBrowseMsg{endpoint: newModel.connectionTextInput.Value()}
				}
		}
	case errMsg:
		m.err = msg
		return m, nil
	}
	m.connectionTextInput, cmd = m.connectionTextInput.Update(msg)
	return m, cmd
}

