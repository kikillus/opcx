// Copyright 2018-2020 opcua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package opcutil

import (
	"context"
	"strconv"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/errors"
	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/ua"
)

type NodeDef struct {
	NodeID      *ua.NodeID
	NodeClass   ua.NodeClass
	BrowseName  string
	Description string
	AccessLevel ua.AccessLevelType
	Path        string
	DataType    string
	Writable    bool
	Unit        string
	Scale       string
	Min         string
	Max         string
}

func (n NodeDef) Records() []string {
	return []string{n.BrowseName, n.DataType, n.NodeID.String(), n.Unit, n.Scale, n.Min, n.Max, strconv.FormatBool(n.Writable), n.Description}
}

func getChildren(ctx context.Context, n *opcua.Node) ([]*opcua.Node, error) {
	refs, err := n.ReferencedNodes(ctx, 33, ua.BrowseDirectionForward, ua.NodeClassAll, true)
	if err != nil {
		return nil, errors.Errorf("References: %s", err)
	}
	var nodes []*opcua.Node
	nodes = append(nodes, refs...)
	return nodes, nil
}

func getChildrenRecursive(ctx context.Context, rootNode *opcua.Node) ([]*opcua.Node, error){
	var collectedNodes []*opcua.Node
	// creating recursive func with collectedNotes in closure (??)
	var recursiveFunc func(ctx context.Context, n *opcua.Node) error
	recursiveFunc = func(ctx context.Context, n *opcua.Node) error {

		// collect chidren
		children, err := getChildren(ctx, n)

		// add current node to slice if it is a leaf
		if len(children) == 0 {
			collectedNodes = append(collectedNodes, n)
		}
		if err != nil {
			return err
		}
		// run recurice function for children
		for _, child := range children {
			if err := recursiveFunc(ctx, child); err != nil{
				return err
			}
		}
		// base case no children
		return nil
	}
	if err := recursiveFunc(ctx, rootNode); err != nil {
		return nil, err
	}
	return collectedNodes, nil
	}

func getNodeAttributes(ctx context.Context, n *opcua.Node) (NodeDef, error) {
	attrs, err := n.Attributes(ctx, ua.AttributeIDNodeClass, ua.AttributeIDBrowseName, ua.AttributeIDDescription, ua.AttributeIDAccessLevel, ua.AttributeIDDataType)
	if err != nil {
		return NodeDef{}, err
	}

	var def = NodeDef{
		NodeID: n.ID,
	}

	switch err := attrs[0].Status; err {
	case ua.StatusOK:
		def.NodeClass = ua.NodeClass(attrs[0].Value.Int())
	default:
		return NodeDef{}, err
	}

	switch err := attrs[1].Status; err {
	case ua.StatusOK:
		def.BrowseName = attrs[1].Value.String()
	default:
		return NodeDef{}, err
	}

	switch err := attrs[2].Status; err {
	case ua.StatusOK:
		def.Description = attrs[2].Value.String()
	case ua.StatusBadAttributeIDInvalid:
		// ignore
	default:
		return NodeDef{}, err
	}

	switch err := attrs[3].Status; err {
	case ua.StatusOK:
		def.AccessLevel = ua.AccessLevelType(attrs[3].Value.Int())
		def.Writable = def.AccessLevel&ua.AccessLevelTypeCurrentWrite == ua.AccessLevelTypeCurrentWrite
	case ua.StatusBadAttributeIDInvalid:
		// ignore
	default:
		return NodeDef{}, err
	}

	switch err := attrs[4].Status; err {
	case ua.StatusOK:
		switch v := attrs[4].Value.NodeID().IntID(); v {
		case id.DateTime:
			def.DataType = "time.Time"
		case id.Boolean:
			def.DataType = "bool"
		case id.SByte:
			def.DataType = "int8"
		case id.Int16:
			def.DataType = "int16"
		case id.Int32:
			def.DataType = "int32"
		case id.Byte:
			def.DataType = "byte"
		case id.UInt16:
			def.DataType = "uint16"
		case id.UInt32:
			def.DataType = "uint32"
		case id.UtcTime:
			def.DataType = "time.Time"
		case id.String:
			def.DataType = "string"
		case id.Float:
			def.DataType = "float32"
		case id.Double:
			def.DataType = "float64"
		default:
			def.DataType = attrs[4].Value.NodeID().String()
		}
	case ua.StatusBadAttributeIDInvalid:
		// ignore
	default:
		return NodeDef{}, err
	}

	return def, nil
}

func ReadValue(ctx context.Context, client *opcua.Client, nodesToRead []*ua.ReadValueID) (interface{}, error) {
	req := &ua.ReadRequest{
		NodesToRead: nodesToRead,
	}
	response, err := client.Read(ctx, req)
	if err != nil {
		return nil, err
	}
	return response.Results[0].Value.Value(), nil
}

// Exported Browse function to be used in the TUI application
func Browse(ctx context.Context, n *opcua.Node) (NodeDef, error) {
	return getNodeAttributes(ctx, n)
}

func GetChildren(ctx context.Context, n *opcua.Node) ([]*opcua.Node, error) {
	return getChildren(ctx, n)
}

func GetChildrenRecursive(ctx context.Context, rootNode *opcua.Node) ([]*opcua.Node, error){
	return getChildrenRecursive(ctx, rootNode)
}