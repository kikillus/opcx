package app

import (
	"opcx/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
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
		return m, nil
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

