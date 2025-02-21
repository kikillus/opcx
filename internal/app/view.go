package app

import (
	"fmt"
	"opcx/internal/ui"
)

func (m model) View() string {
	content, header, footer := ui.RenderView(m.state, m.nav, m.activeNode, m.client.ReadNodeValue, m.connectionTextInput)

	var rendered string
	switch m.state {
	case ui.ViewStateBrowse, ui.ViewStateRecursive:
		m.viewport.SetContent(content)
		rendered = fmt.Sprintf("%s\n%s\n%s", header, m.viewport.View(), footer)
	default:
		rendered = fmt.Sprintf("%s\n%s\n%s\n", header, content, footer)
	}
	return rendered
}
