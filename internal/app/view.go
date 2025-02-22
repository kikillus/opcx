package app

import (
	"fmt"
	"opcx/internal/ui"

	"github.com/charmbracelet/bubbles/textinput"
)

func (m model) View() string {
	var nav *ui.Navigation
	var connectionTextInput textinput.Model
	switch m.state{
	case ui.ViewStateBrowse:
		nav = m.browseView.nav
	case ui.ViewStateRecursive:
		nav = m.recursiveView.nav
	case ui.ViewStateMonitor:
		nav = m.monitorView.nav
	case ui.ViewStateConnection:
		connectionTextInput = m.connectionView.connectionTextInput
	}
	content, header, footer := ui.RenderView(m.state, nav, m.detailsView.activeNode, m.client.ReadNodeValue, connectionTextInput, m.monitoredNodes)

	var rendered string
	switch m.state {
	case ui.ViewStateBrowse:
		m.browseView.viewport.SetContent(content)
		rendered = fmt.Sprintf("%s\n%s\n%s", header, m.browseView.viewport.View(), footer)
	case ui.ViewStateRecursive:
		m.recursiveView.viewport.SetContent(content)
		rendered = fmt.Sprintf("%s\n%s\n%s", header, m.recursiveView.viewport.View(), footer)
	default:
		rendered = fmt.Sprintf("%s\n%s\n%s\n", header, content, footer)
	}
	return rendered
}
