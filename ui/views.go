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

func RenderView(state ViewState, nav *Navigation, activeNode opcutil.NodeDef, readNodeValue func(*ua.NodeID) string, connectionTextInput textinput.Model) (string, string,  string){
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
		return "Unkown state", "", ""
	}
}

func renderRecursiveView(nav *Navigation) (string, string, string){
	header := fmt.Sprintf("All leaf children of: %s\n\n", nav.CurrentNode().BrowseName)
	s := ""
	for i, node := range nav.CurrentNodes{
		cursor := " "
		if nav.Cursor == i {
			cursor = ">"}
		s += fmt.Sprintf("%s BrowseName: %s - NodeID: %s - DataType: %s\n", cursor,node.BrowseName, node.NodeID, node.DataType)
	}
	footer := "\n[q]uit - toggle [c]hildren"
	return s, header, footer
}

func renderDetailView(node opcutil.NodeDef, readNodeValue func(*ua.NodeID) string) (string, string, string){
	header := "OPC UA Node Detail\n\n"
	s := fmt.Sprintf("BrowseName: %s\n", node.BrowseName)
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
	footer := "\n[q]uit - toggle [v]iew"
	return s, header, footer
}

func renderBrowseView(nav *Navigation) (string,string, string){
	header := "OPC UA Node Browser"

	content := ""
	for i, node := range nav.CurrentNodes {
		cursor := " "
		if nav.Cursor == i {
			cursor = ">"
		}

		if node.DataType == "" {
		content += fmt.Sprintf("%s %s\n", cursor, node.BrowseName)
		} else {
			content += fmt.Sprintf("%s %s (%s)\n", cursor, node.BrowseName, node.DataType)
		}
	}
	footer := ""
	path := buildPath(nav.Path)
	if path != "" {
		footer += "\n" + fmt.Sprintf("Path: %s", path)
	}
	footer += "\n[q]uit - toggle [v]iew - toogle leaf [c]hildren"
	return content, header, footer
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

func renderConnectionView(connectionTextInput textinput.Model) (string, string, string) {
	s := "Connect to OPC Server\n"
	s += connectionTextInput.View()
	return s, "", ""
}