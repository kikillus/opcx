package app

import (
	"opcx/internal/opc"
	"opcx/internal/ui"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func NewRecursiveViewModel() RecursiveViewModel {
	vp := viewport.New(0, 0)
	vp.YPosition = 3

	return RecursiveViewModel{
		nav: ui.NewNavigation(),
		viewport: &vp,
	}
}

type TransitionRecursiveToBrowseMsg struct {
}

func (m RecursiveViewModel) Update(msg tea.Msg) (RecursiveViewModel, tea.Cmd) {
	if m.nav.Cursor < 0 {
		m.nav.Cursor = 0
	} else if m.nav.Cursor >= len(m.nav.CurrentNodes) {
		m.nav.Cursor = len(m.nav.CurrentNodes) - 1
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up", "k":
			if m.nav.Cursor > 0 {
				newModel := m
				newModel.nav.Cursor--
				if newModel.nav.Cursor < newModel.viewport.YOffset {
					newModel.viewport.SetYOffset(newModel.nav.Cursor)
				}
				return newModel, nil
			}
			return m, nil
		case "down", "j":
			if m.nav.Cursor < len(m.nav.CurrentNodes)-1 {
				newModel := m
				newModel.nav.Cursor++
				// Adjust viewport if cursor is beyond visible area
				if newModel.nav.Cursor >= newModel.viewport.YOffset+newModel.viewport.Height {
					newModel.viewport.SetYOffset(newModel.nav.Cursor - newModel.viewport.Height + 1)
				}
				return newModel, nil
			}
			return m, nil
		case "c":
			return m, func () tea.Msg {
				return TransitionRecursiveToBrowseMsg{}
			}
		}
	case []opc.NodeDef:
		newModel := m
		newModel.nav.CurrentNodes = msg
		newModel.nav.Cursor = 0
		newModel.viewport.SetYOffset(0) // Reset viewport position when loading new nodes
		return newModel, nil
	}
	return m, nil
}
