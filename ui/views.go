package ui

import (
	"fmt"
	opcutil "opc-tui/opc/util"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/gopcua/opcua/ua"
)

type ViewState int

const (
	ViewStateBrowse ViewState = iota
	ViewStateDetail
	ViewStateConnection
	ViewStateRecursive
)

func RenderView(state ViewState, nav *Navigation, activeNode opcutil.NodeDef, readNodeValue func(*ua.NodeID) string, connectionTextInput textinput.Model) string {
	switch state {
	case ViewStateBrowse:
		return renderBrowseView(nav)
	case ViewStateDetail:
		return renderDetailView(activeNode, readNodeValue)
	case ViewStateConnection:
		return renderConnectionView(connectionTextInput)
	case ViewStateRecursive:
		return renderRecursiveView(nav)
	default:
		return "Unkown state"
	}
}

func renderRecursiveView(nav *Navigation) string{
	s := fmt.Sprintf("All leaf children of: %s\n\n", nav.CurrentNode().BrowseName)
	for _, node := range nav.CurrentNodes{
		s += fmt.Sprintf("BrowseName: %s - NodeID: %s - DataType: %s\n", node.BrowseName, node.NodeID, node.DataType)
	}
	return s
}

func renderDetailView(node opcutil.NodeDef, readNodeValue func(*ua.NodeID) string) string {
	s := "OPC UA Node Detail\n\n"
	s += fmt.Sprintf("BrowseName: %s\n", node.BrowseName)
	s += fmt.Sprintf("NodeID: %s\n", node.NodeID)
	s += fmt.Sprintf("Description: %s\n", node.Description)
	s += fmt.Sprintf("AccessLevel: %s\n", node.AccessLevel)
	s += fmt.Sprintf("Path: %s\n", node.Path)
	s += fmt.Sprintf("DataType: %s\n", node.DataType)
	s += fmt.Sprintf("Writable: %t\n", node.Writable)
	s += fmt.Sprintf("Unit: %s\n", node.Unit)
	s += fmt.Sprintf("Scale: %s\n", node.Scale)
	s += fmt.Sprintf("Min: %s\n", node.Min)
	s += fmt.Sprintf("Max: %s\n", node.Max)
	value := readNodeValue(node.NodeID)
	if !(value == "default") {
		s += fmt.Sprintf("Value: %s\n", value)
	}
	s += "\n[q]uit - toggle [v]iew\n"
	return s
}

func renderBrowseView(nav *Navigation) string {
	s := "OPC UA Node Browser\n\n"

	for i, node := range nav.CurrentNodes {
		cursor := " "
		if nav.Cursor == i {
			cursor = ">"
		}

		if node.DataType == "" {
			s += fmt.Sprintf("%s %s\n", cursor, node.BrowseName)
		} else {
			s += fmt.Sprintf("%s %s (%s)\n", cursor, node.BrowseName, node.DataType)
		}
	}
	path := buildPath(nav.Path)
	if path != "" {
		s += "\n" + fmt.Sprintf("Path: %s", path)
	}
	s += "\n[q]uit - toggle [v]iew\n"
	return s
}

func buildPath(path []opcutil.NodeDef) string {
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

func renderConnectionView(connectionTextInput textinput.Model) string {
	s := "Connect to OPC Server\n"
	s += connectionTextInput.View()
	return s
}