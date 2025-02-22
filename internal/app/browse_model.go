package app

import (
	"opcx/internal/opc"
	"opcx/internal/ui"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func NewBrowseViewModel() BrowseViewModel {
	vp := viewport.New(0, 0)
	vp.YPosition = 3

	return BrowseViewModel{
		viewport: &vp,
		nav:      ui.NewNavigation(),
	}
}

type FetchChildrenMsg struct {
	NodeID string
}

type FetchChildrenRecursiveMsg struct {
	NodeID string
}

type ChangeViewStateMsg struct {
	NewState ui.ViewState
}

func (m BrowseViewModel) Update(msg tea.Msg) (BrowseViewModel, tea.Cmd) {
	if m.nav.Cursor < 0 {
		m.nav.Cursor = 0
	} else if m.nav.Cursor >= len(m.nav.CurrentNodes) {
		m.nav.Cursor = len(m.nav.CurrentNodes) - 1
	}
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
		case "enter", "right", "l":
			if len(m.nav.CurrentNodes) == 0 {
				return m, nil
			}
			newModel := m
			currentNode := m.nav.CurrentNodes[m.nav.Cursor]
			newModel.nav.Forward(m.nav.CurrentNodes[m.nav.Cursor])
			return newModel, func() tea.Msg { return FetchChildrenMsg{NodeID: currentNode.NodeID.String()} } // FIXME
		case "left", "h":
			newModel := m
			if newModel.nav.Backward() {
				newCurrentNode := m.nav.CurrentNode()
				newModel := m
				return newModel, func() tea.Msg { return FetchChildrenMsg{NodeID: newCurrentNode.NodeID.String()} } // FIXME
			}
			return m, nil
		case "v":
			newModel := m
			newModel.activeNode = m.nav.CurrentNodes[m.nav.Cursor]
			return newModel, func() tea.Msg { return ChangeViewStateMsg{NewState: ui.ViewStateDetail} }
		case "c":
			newModel := m
			currentNode := m.nav.CurrentNodes[m.nav.Cursor]
			newModel.activeNode = currentNode
			newModel.nav.Forward(m.nav.CurrentNodes[m.nav.Cursor])
			return newModel, tea.Batch(
				func() tea.Msg {
					return FetchChildrenRecursiveMsg{NodeID: currentNode.NodeID.String()}
				},
				func() tea.Msg {
					return ChangeViewStateMsg{NewState: ui.ViewStateRecursive}
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
	case errMsg:
		m.err = msg
		return m, nil
	}
	return m, cmd
}
