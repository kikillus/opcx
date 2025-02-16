package main

import (
	"fmt"
	"log"

	opcservice "opc-tui/opc/service"
	opcutil "opc-tui/opc/util"
	"opc-tui/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gopcua/opcua/ua"
)

type model struct {
	nav *ui.Navigation
	service *opcservice.Service
	active_node opcutil.NodeDef
	state 	ui.ViewState
	width int
	height int
}

func initialModel(service *opcservice.Service) model {
	return model{nav: ui.NewNavigation(),service: service, state: ui.ViewStateBrowse}
}

func (m model) Init() tea.Cmd {
	return m.fetchChildren(m.nav.CurrentNode().NodeID)
}

func (m model) fetchChildren(nodeID *ua.NodeID) tea.Cmd {
	return func() tea.Msg {
		children, err := m.service.GetChildren(nodeID)
		if err != nil {
			log.Fatalf("browse error: %s", err)
		}
		return children
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.nav.Cursor < 0 {
		m.nav.Cursor = 0
	} else if m.nav.Cursor >= len(m.nav.CurrentNodes) {
		m.nav.Cursor = len(m.nav.CurrentNodes) - 1
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.nav.Cursor > 0 {
				m.nav.Cursor--
			}
		case "down", "j":
			if m.nav.Cursor < len(m.nav.CurrentNodes)-1 {
				m.nav.Cursor++
			}
		case "enter", "right", "l":
			if len(m.nav.CurrentNodes) == 0 {
				return m, nil
			}
			newModel := m
			newModel.nav.Forward(m.nav.CurrentNodes[m.nav.Cursor])
			children := m.fetchChildren(m.nav.CurrentNodes[m.nav.Cursor].NodeID)
			return newModel, children
		case "left", "h":
			newModel := m
			if newModel.nav.Backward() {
				newCurrentNode := m.nav.CurrentNode()
				newModel := m
				children := m.fetchChildren(newCurrentNode.NodeID)
				return newModel, children
			}
			return m, nil
		case "v":
			if m.state == ui.ViewStateBrowse {
				newModel := m
				newModel.state = ui.ViewStateDetail
				newModel.active_node = m.nav.CurrentNodes[m.nav.Cursor]
				return newModel, nil
			}
			newModel := m
			newModel.state = ui.ViewStateBrowse
			return newModel, nil
		}
	case tea.WindowSizeMsg:
		newModel := m
		newModel.width = msg.Width
		newModel.height = msg.Height
		return newModel, nil
	case []opcutil.NodeDef:
		if len(msg) == 0 {
			newModel := m
			newModel.nav.Path= m.nav.Path[:len(m.nav.Path)-1]
			return newModel, nil
		}
		newModel := m
		newModel.nav.CurrentNodes = msg
		newModel.nav.Cursor = 0
		return newModel, nil
	}
	return m, nil
}

func (m model) View() string {
	return ui.RenderView(m.state, m.nav, m.active_node, m.readNodeValue)}

func (m model) readNodeValue(nodeID *ua.NodeID) (string) {
	fmt.Println("got here")
	value, err := m.service.ReadNodeValue(nodeID)
	if err != nil {
		log.Fatalf("read error: %s", err)
	}
	return value
}

func main() {
	endpoint := "opc.tcp://localhost:4840"
	service, err := opcservice.NewService(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close()

	p := tea.NewProgram(initialModel(service), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
