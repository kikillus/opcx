package opcservice

import (
	"context"
	opcutil "opc-tui/opc/util"
	"strconv"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

type Service struct {
	client *opcua.Client
	ctx    context.Context
}

func (s *Service) GetNodeFromID(nodeID *ua.NodeID) *opcua.Node {
	return s.client.Node(nodeID)
}

func NewService(endpoint string) (*Service, error) {
	client, err := opcua.NewClient(endpoint)
	if err != nil {
		return nil, err
	}
	if err := client.Connect(context.Background()); err != nil {
		return nil, err
	}
	ctx := context.Background()
	return &Service{client: client, ctx: ctx}, nil
}

func (s *Service) Close() error {
	return s.client.Close(s.ctx)
}

func (s *Service) GetChildren(nodeID *ua.NodeID) ([]opcutil.NodeDef, error) {
	node := s.client.Node(nodeID)

	nodes, err := opcutil.GetChildren(s.ctx, node)
	if err != nil {
		return nil, err
	}

	nodeDefs := make([]opcutil.NodeDef, 0, len(nodes))
	for _, n := range nodes {
		def, err := opcutil.Browse(s.ctx, n)
		if err != nil {
			return nil, err
		}
		nodeDefs = append(nodeDefs, def)
	}
	return nodeDefs, nil
}

func (s *Service) GetChildrenRecursive(rootNodeID *ua.NodeID) ([]opcutil.NodeDef, error) {
	rootNode := s.client.Node(rootNodeID)

	nodes, err := opcutil.GetChildrenRecursive(s.ctx, rootNode)
	if err != nil {
		return nil, err
	}

	nodeDefs := make([]opcutil.NodeDef, 0, len(nodes))
	for _, n := range nodes {
		def, err := opcutil.Browse(s.ctx, n)
		if err != nil {
			return nil,err
		}
		nodeDefs = append(nodeDefs, def)
	}
	return nodeDefs, nil
}

func (s *Service) ReadNodeValue(nodeID *ua.NodeID) (string, error) {
	nodesToRead := []*ua.ReadValueID{
		{NodeID: nodeID},
	}
	value, err := opcutil.ReadValue(s.ctx, s.client, nodesToRead)
	if err != nil {
		return "", err
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(v), nil
	case time.Time:
		return v.String(), nil
	default:
		return "default", nil
	}
}
