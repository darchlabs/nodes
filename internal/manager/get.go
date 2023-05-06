package manager

import (
	"sort"

	"github.com/pkg/errors"
)

var ErrNodeNotFound = errors.New("manager: node not found")

func (m *Manager) Get(id string) (*NodeInstance, error) {
	nodeInstance, ok := m.nodes[id]
	if !ok {
		return nil, ErrNodeNotFound
	}

	return nodeInstance, nil
}

func (m *Manager) GetAll() []*NodeInstance {
	nodes := make([]*NodeInstance, 0)

	for _, nodeCommand := range m.nodes {
		nodes = append(nodes, nodeCommand)
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Config.CreatedAt.Before(nodes[j].Config.CreatedAt)
	})

	return nodes
}
