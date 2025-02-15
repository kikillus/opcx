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

type viewState int

const (
	viewStateBrowse viewState = iota
	viewStateDetail
)
type model struct {
	nav *ui.Navigation
	service *opcservice.Service
	active_node opcutil.NodeDef
	state 	viewState
	width int
	height int
}

func initialModel(service *opcservice.Service) model {
	return model{nav: ui.NewNavigation(),service: service, state: viewStateBrowse}
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
			if m.state == viewStateDetail {
				newModel := m
				newModel.state = viewStateBrowse
				return newModel, nil
			}
			newModel := m
			newModel.active_node = m.nav.CurrentNodes[m.nav.Cursor]
			newModel.state = viewStateDetail
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
	switch m.state {
	case viewStateBrowse:
		return m.viewBrowse()
	case viewStateDetail:
		return m.viewDetail()
	default:
		return "Unkown state"
	}
}
func (m model) viewDetail() string {
	s := "OPC UA Node Detail\n\n"
	s += fmt.Sprintf("BrowseName: %s\n", m.active_node.BrowseName)
	s += fmt.Sprintf("NodeID: %s\n", m.active_node.NodeID)
	s += fmt.Sprintf("Description: %s\n", m.active_node.Description)
	s += fmt.Sprintf("AccessLevel: %s\n", m.active_node.AccessLevel)
	s += fmt.Sprintf("Path: %s\n", m.active_node.Path)
	s += fmt.Sprintf("DataType: %s\n", m.active_node.DataType)
	s += fmt.Sprintf("Writable: %t\n", m.active_node.Writable)
	s += fmt.Sprintf("Unit: %s\n", m.active_node.Unit)
	s += fmt.Sprintf("Scale: %s\n", m.active_node.Scale)
	s += fmt.Sprintf("Min: %s\n", m.active_node.Min)
	s += fmt.Sprintf("Max: %s\n", m.active_node.Max)
	value:= m.readNodeValue(m.active_node.NodeID)
	if !(value == "default") {
		s += fmt.Sprintf("Value: %s\n", value)
	}
	s += "\n[q]uit - toggle [v]iew\n"
	return s
}
func (m model) viewBrowse() string {
	s := "OPC UA Node Browser\n\n"
	for i, node := range m.nav.CurrentNodes{
		cursor := " "
		if m.nav.Cursor == i {
			cursor = ">"
		}
		if node.DataType == "" {
		s += fmt.Sprintf("%s %s\n", cursor, node.BrowseName)} else {
		s += fmt.Sprintf("%s %s (%s)\n", cursor, node.BrowseName, node.DataType)
		}
	}
	path := ""
	for _, parent := range m.nav.Path {
		if parent.BrowseName == "" {
			continue
		}
		if path == "" {
			path += parent.BrowseName
		} else {
			path += " > " + parent.BrowseName
		}
	}
	if !(path == "") {
	s += "\n" + fmt.Sprintf("Path: %s", path)
	}
	s += "\n[q]uit - toggle [v]iew\n"
	return s
}

func (m model) readNodeValue(nodeID *ua.NodeID) (string) {
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
