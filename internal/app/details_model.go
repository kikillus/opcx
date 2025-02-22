package app

import (
	"opcx/internal/opc"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type DetailsViewModel struct {
	viewport   *viewport.Model
	activeNode opc.NodeDef
}

func NewDetailsViewModel() DetailsViewModel {
	vp := viewport.New(0, 0)
	vp.YPosition = 3

	return DetailsViewModel{
		viewport: &vp,
	}
}

type TransitionDetailToBrowseMsg struct {
}
func (m DetailsViewModel) Update(msg tea.Msg) (DetailsViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "v":
			return m, func() tea.Msg {
				return TransitionDetailToBrowseMsg{}
			}

		}
	}
	return m, nil
}

