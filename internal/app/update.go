package app

import (
	"opcx/internal/opc"
	"opcx/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		newModel := m
		newModel.width = msg.Width
		newModel.height = msg.Height

		headerHeight := 6
		footerHeight := 0
		verticalMarginHeight := headerHeight + footerHeight

		newModel.viewport.Width = msg.Width
		newModel.viewport.Height = msg.Height - verticalMarginHeight
		return newModel, nil
	}
	var newModel tea.Model
	var cmd tea.Cmd
	switch m.state {
	case ui.ViewStateConnection:
		newModel, cmd = m.updateConnectionView(msg)
	case ui.ViewStateDetail:
		newModel, cmd = m.updateDetailView(msg)
	case ui.ViewStateBrowse:
		newModel, cmd = m.updateBrowseView(msg)
	case ui.ViewStateRecursive:
		newModel, cmd = m.updateRecursiveView(msg)
	default:
		newModel = m
		cmd = nil
	}
	return newModel, cmd
}

func (m model) updateBrowseView(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			newModel.nav.Forward(m.nav.CurrentNodes[m.nav.Cursor])
			children := m.client.FetchChildren(m.nav.CurrentNodes[m.nav.Cursor].NodeID)
			return newModel, children
		case "left", "h":
			newModel := m
			if newModel.nav.Backward() {
				newCurrentNode := m.nav.CurrentNode()
				newModel := m
				children := m.client.FetchChildren(newCurrentNode.NodeID)
				return newModel, children
			}
			return m, nil
		case "v":
			newModel := m
			newModel.state = ui.ViewStateDetail
			newModel.activeNode = m.nav.CurrentNodes[m.nav.Cursor]
			return newModel, nil
		case "c":
			newModel := m
			currentNode := m.nav.CurrentNodes[m.nav.Cursor]
			newModel.activeNode = currentNode
			newModel.nav.Forward(m.nav.CurrentNodes[m.nav.Cursor])
			newModel.state = ui.ViewStateRecursive
			recursiveChildren := m.client.FetchChildrenRecursive(currentNode.NodeID)
			return newModel, recursiveChildren
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

func (m model) updateRecursiveView(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			children := m.client.FetchChildren(newCurrentNode.NodeID)
			newModel.state = ui.ViewStateBrowse
			return newModel, children
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

func (m model) updateDetailView(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			newModel := m
			newModel.state = ui.ViewStateBrowse
			return newModel, nil
		}
	}
	return m, nil
}

func (m model) updateConnectionView(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			client, err := opc.NewClient(newModel.connectionTextInput.Value())
			if err != nil {
				return newModel, func() tea.Msg { return errMsg(err) }
			}
			newModel.client = client

			cmd = newModel.client.FetchChildren(newModel.nav.CurrentNode().NodeID)
			newModel.state = ui.ViewStateBrowse
			return newModel, cmd
		}
	case errMsg:
		m.err = msg
		return m, nil
	}
	m.connectionTextInput, cmd = m.connectionTextInput.Update(msg)
	return m, cmd
}
