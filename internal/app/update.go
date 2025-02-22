package app

import (
	"opcx/internal/opc"
	"opcx/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gopcua/opcua/ua"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 6
		footerHeight := 0
		verticalMarginHeight := headerHeight + footerHeight

		if m.browseView.viewport != nil {
			m.browseView.viewport.Width = msg.Width
			m.browseView.viewport.Height = msg.Height - verticalMarginHeight
		}
		if m.detailsView.viewport != nil {
			m.detailsView.viewport.Width = msg.Width
			m.detailsView.viewport.Height = msg.Height - verticalMarginHeight
		}
		if m.recursiveView.viewport != nil {
			m.recursiveView.viewport.Width = msg.Width
			m.recursiveView.viewport.Height = msg.Height - verticalMarginHeight
		}
		return m, nil
	case TransitionConnectToBrowseMsg:
		var err error
		m.client, err = opc.NewClient(msg.endpoint)
		if err != nil {
			return m, func() tea.Msg { return errMsg(err) }
		}
		m.state = ui.ViewStateBrowse
		children := m.client.FetchChildren(ua.NewNumericNodeID(0, 84))
		return m, children
	case TransitionBrowseToDetailMsg:
		m.state = ui.ViewStateDetail
		m.detailsView.activeNode = opc.NodeDef{NodeID: msg.nodeID}
		return m, nil
	case TransitionDetailToBrowseMsg:
		m.state = ui.ViewStateBrowse
		return m, nil
	case TransitionBrowseToRecursiveMsg:
		m.state = ui.ViewStateRecursive
		m.recursiveView.nav.ActiveNode = msg.rootNode
		children := m.client.FetchChildrenRecursive(msg.rootNode.NodeID)
		return m, children
	case TransitionRecursiveToBrowseMsg:
		m.state = ui.ViewStateBrowse
		return m, nil
	case FetchChildrenMsg:
		children := m.client.FetchChildren(msg.NodeID)
		return m, children
	}
	var newModel model
	var cmd tea.Cmd

	switch m.state {
	case ui.ViewStateConnection:
		newModel = m
		newModel.connectionView, cmd = m.connectionView.Update(msg)
	case ui.ViewStateDetail:
		newModel = m
		newModel.detailsView, cmd = m.detailsView.Update(msg)
	case ui.ViewStateBrowse:
		newModel = m
		newModel.browseView, cmd = m.browseView.Update(msg)
	case ui.ViewStateRecursive:
		newModel = m
		newModel.recursiveView, cmd = m.recursiveView.Update(msg)
	default:
		newModel = m
		cmd = nil
	}
	return newModel, cmd
}

// TODO Define state transitions
