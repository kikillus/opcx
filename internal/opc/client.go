package opc

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gopcua/opcua/ua"
)

type Client struct {
	service *Service
}

type errMsg error

func NewClient(endpoint string) (*Client, error) {
	service, err := NewService(endpoint)
	if err != nil {
		return nil, err
	}
	return &Client{service: service}, nil
}

func (c *Client) ReadNodeValue(nodeID *ua.NodeID) (string, error) {
	return c.service.ReadNodeValue(nodeID)
}

func (c *Client) FetchChildren(nodeID *ua.NodeID) tea.Cmd {
	return func() tea.Msg {
		children, err := c.service.GetChildren(nodeID)
		if err != nil {
			return errMsg(err)
		}
		return children
	}
}

func (c *Client) FetchChildrenRecursive(rootNode *ua.NodeID) tea.Cmd {
	return func() tea.Msg {
		children, err := c.service.GetChildrenRecursive(rootNode)
		if err != nil {
			return errMsg(err)
		}
		return children
	}
}
