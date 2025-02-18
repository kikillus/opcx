package main

import (
	"fmt"
	"log"

	opcservice "opc-tui/opc/service"
	opcutil "opc-tui/opc/util"
	"opc-tui/ui"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
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
	connectionTextInput textinput.Model
	err error
	viewport *viewport.Model  // Change to pointer
}

type (
	errMsg error
)

func initialModel() model {
	connectionText := textinput.New()
	connectionText.Placeholder = "opc.tcp://127.0.0.1:4840"
	connectionText.SetValue("opc.tcp://127.0.0.1:4840")
	connectionText.CharLimit = 128
	connectionText.Width = 50
	connectionText.Focus()
	vp := viewport.New(0,0)
	vp.YPosition = 3 // Add this to position viewport below header
	return model{nav: ui.NewNavigation(), viewport: &vp, connectionTextInput: connectionText, state: ui.ViewStateConnection}
}


func (m model) Init() tea.Cmd {
	return textinput.Blink
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type){
	case tea.WindowSizeMsg:
		newModel := m
		newModel.width = msg.Width
		newModel.height = msg.Height

		headerHeight := 6
		footerHeight := 0
		verticalMarginHeight := headerHeight + footerHeight

		newModel.viewport.Width = msg.Width
		newModel.viewport.Height = msg.Height - verticalMarginHeight
		return newModel,nil
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

func (m model) updateBrowseView(msg tea.Msg) (tea.Model, tea.Cmd){
	if m.nav.Cursor < 0 {
		m.nav.Cursor = 0
	} else if m.nav.Cursor >= len(m.nav.CurrentNodes) {
		m.nav.Cursor = len(m.nav.CurrentNodes) - 1
	}
	var cmd tea.Cmd
	switch msg := msg.(type){
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
		switch msg.String(){
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
				newModel := m
				newModel.state = ui.ViewStateDetail
				newModel.active_node = m.nav.CurrentNodes[m.nav.Cursor]
				return newModel, nil
		case "c":
			newModel := m
			currentNode := m.nav.CurrentNodes[m.nav.Cursor]
			newModel.active_node = currentNode
			newModel.nav.Forward(m.nav.CurrentNodes[m.nav.Cursor])
			newModel.state = ui.ViewStateRecursive
			recursiveChildren := m.fetchChildrenRecursive(currentNode.NodeID)
			return newModel, recursiveChildren
		}
	case []opcutil.NodeDef:
		newModel := m
		if len(msg) == 0 {
			newModel.nav.Path= m.nav.Path[:len(m.nav.Path)-1]
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
            children := m.fetchChildren(newCurrentNode.NodeID)
            newModel.state = ui.ViewStateBrowse
            return newModel, children
        }
    case []opcutil.NodeDef:
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

func (m model) updateDetailView(msg tea.Msg) (tea.Model, tea.Cmd){
	switch msg := msg.(type){
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
		switch msg.String(){
		case "q":
			return m, tea.Quit
		case "v":
			newModel :=m
			newModel.state = ui.ViewStateBrowse
			return newModel, nil
		}
	}
	return m, nil
}

func (m model) updateConnectionView(msg tea.Msg) (tea.Model, tea.Cmd){
	var cmd tea.Cmd
	switch msg := msg.(type){
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
		switch msg.String(){
		case "q":
			return m, tea.Quit
		case "enter":
			newModel := m
			newModel = newModel.connectService(newModel.connectionTextInput.Value())

			cmd  = newModel.fetchChildren(newModel.nav.CurrentNode().NodeID)
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

func (m model) View() string {
	content, header, footer := ui.RenderView(m.state, m.nav, m.active_node, m.readNodeValue, m.connectionTextInput)

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

func (m model) readNodeValue(nodeID *ua.NodeID) (string) {
	value, err := m.service.ReadNodeValue(nodeID)
	if err != nil {
		log.Fatalf("read error: %s", err)
	}
	return value
}

func (m model) connectService(endpoint string) model{
	service, err := opcservice.NewService(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	m.service = service
	return m
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
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

func (m model) fetchChildrenRecursive(rootNode *ua.NodeID) tea.Cmd {
	return func() tea.Msg {
		children, err := m.service.GetChildrenRecursive(rootNode)
		if err != nil {
			log.Fatalf("browse recursive error: %s", err)
		}
		return children
	}
}