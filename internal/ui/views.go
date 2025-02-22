package ui

import (
	"fmt"
	"opcx/internal/opc"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/gopcua/opcua/ua"
)

type ViewState int

const (
	ViewStateBrowse ViewState = iota
	ViewStateDetail
	ViewStateConnection
	ViewStateRecursive
	ViewStateMonitor
)

func RenderView(state ViewState, nav *Navigation, activeNode opc.NodeDef, readNodeValue func(*ua.NodeID) (string, error), connectionTextInput textinput.Model, monitoredNodes []opc.NodeDef) (string, string,  string){
	switch state {
	case ViewStateBrowse:
		return renderBrowseView(nav, monitoredNodes)
	case ViewStateDetail:
		return renderDetailView(activeNode, readNodeValue)
	case ViewStateConnection:
		return renderConnectionView(connectionTextInput)
	case ViewStateRecursive:
		return renderRecursiveView(nav)
	case ViewStateMonitor:
		return renderMonitorView(nav, readNodeValue)
	default:
		return "Unkown state", "", ""
	}
}

func renderRecursiveView(nav *Navigation) (string, string, string){
	header := HeaderStyle.Render(fmt.Sprintf("All leaf children of: %s\n\n", nav.ActiveNode.BrowseName))
	s := ""
	for i, node := range nav.CurrentNodes{
		cursor := " "
		if nav.Cursor == i {
			cursor = ">"}
		s += fmt.Sprintf("%s BrowseName: %s - NodeID: %s - DataType: %s\n", cursor,node.BrowseName, node.NodeID, node.DataType)
	}
	footer := FooterStyle.Render("\n[q]uit - toggle [c]hildren")
	return s, header, footer
}

func renderDetailView(node opc.NodeDef, readNodeValue func(*ua.NodeID) (string, error)) (string, string, string){
	header := HeaderStyle.Render("OPC UA Node Detail")
    labelStyle := lipgloss.NewStyle().Foreground(subtle)
    s := fmt.Sprintf("%s %s\n",
        labelStyle.Render("BrowseName:"),
        node.BrowseName)
	s += fmt.Sprintf("%s %s\n",
		labelStyle.Render("NodeID:"),
		node.NodeID)
	s += fmt.Sprintf("%s %s\n",
		labelStyle.Render("Description:"),
		node.Description)
	s += fmt.Sprintf("%s %s\n",
		labelStyle.Render("AccessLevel:"),
		node.AccessLevel)
	s += fmt.Sprintf("%s %s\n",
		labelStyle.Render("Path:"),
		node.Path)
	s += fmt.Sprintf("%s %s\n",
		labelStyle.Render("DataType:"),
		node.DataType)
	s += fmt.Sprintf("%s %t\n",
		labelStyle.Render("Writable:"),
		node.Writable)
	s += fmt.Sprintf("%s %s\n",
		labelStyle.Render("Unit:"),
		node.Unit)
	s += fmt.Sprintf("%s %s\n",
		labelStyle.Render("Scale:"),
		node.Scale)
	s += fmt.Sprintf("%s %s\n",
		labelStyle.Render("Min:"),
		node.Min)
	s += fmt.Sprintf("%s %s\n",
		labelStyle.Render("Max:"),
		node.Max)
	value, err := readNodeValue(node.NodeID)
	if err != nil {
		s += fmt.Sprintf("Error reading value: %s\n", err)
	}
	if !(value == "default") && value != "" {
		s += fmt.Sprintf("Value: %s\n", value)
	}
	footer := FooterStyle.Render("[q]uit - toggle [v]iew")
	return s, header, footer
}

func renderBrowseView(nav *Navigation, monitoredNodes []opc.NodeDef) (string, string, string) {
    header := HeaderStyle.Render("OPC UA Node Browser")
    content := ""
    for i, node := range nav.CurrentNodes {
		var prefix string
		if node.IsInSlice(monitoredNodes) {
			prefix = "[m] "
		} else {
			prefix = "[ ] "
		}
        if i == nav.Cursor {
            content += SelectedStyle.Render(prefix + "▸ " + node.BrowseName) + "\n"
        } else {
            content += prefix + "  " + node.BrowseName + "\n"
        }
    }
	path := buildPath(nav.Path)
	footer :=""
	if path != "" {
		footer += "\n" + fmt.Sprintf("Path: %s", path)
	}
	footer += "\n[q]uit - toggle [v]iew - toogle leaf [c]hildren - toggle [m]onitor view - [space] add to monitor"
	return content, header, FooterStyle.Render(footer)
}

func buildPath(path []opc.NodeDef) string {
	path_as_string := ""
	for _, parent := range path {
		if parent.BrowseName == "" {
			continue
		}
		if path_as_string == "" {
			path_as_string += parent.BrowseName
		} else {
			path_as_string += " > " + parent.BrowseName
		}

	}
	return path_as_string
}

func renderConnectionView(connectionTextInput textinput.Model) (string, string, string) {
    header := HeaderStyle.Render("OPC UA Connection")
    inputBox := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1).
        Render(connectionTextInput.View())
    footer := FooterStyle.Render("Enter to connect - Ctrl+c to quit")
    return inputBox, header, footer
}

func renderMonitorView(nav *Navigation, readNodeValue func(*ua.NodeID) (string, error))(string, string, string) {
	header := HeaderStyle.Render("OPC UA Monitor")
	content := ""
	for i, node := range nav.CurrentNodes {
		value, err := readNodeValue(node.NodeID)
		if err != nil {
			value = "Error reading value: " + err.Error()
		}
		if i == nav.Cursor {
			content += SelectedStyle.Render("▸ " + node.BrowseName + " - Value: " + value + " - ID: " + node.NodeID.String()) + "\n"
		} else {
			content += "  " + node.BrowseName + " - Value " + value + " - ID: " + node.NodeID.String()+  "\n"
		}
	}
	footer := FooterStyle.Render("\n[q]uit - toggle [m]onitor view")
	return content, header, footer
}
