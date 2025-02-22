package app

import (
	"opcx/internal/opc"
	"opcx/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

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
			newModel := m
			newModel.nav.Backward()
			newCurrentNode := m.nav.CurrentNode()
			return newModel, tea.Batch(
				func() tea.Msg {
					return ChangeViewStateMsg{NewState: ui.ViewStateBrowse}
				},
				func() tea.Msg {
					return FetchChildrenMsg{NodeID: newCurrentNode.NodeID.String()}
				},
			)
		}
	case []opc.NodeDef:
		newModel := m
		if len(msg) == 0 {
			newModel.nav.Path = m.nav.Path[:len(m.nav.Path)-1]
			newModel.viewport.SetYOffset(0) // Reset viewport position
			return newModel, nil
		}
		newModel.nav.CurrentNodes = msg
		newModel.nav.Cursor = 0
		newModel.viewport.SetYOffset(0) // Reset viewport position when loading new nodes
		return newModel, nil
	}
	return m, nil
}
