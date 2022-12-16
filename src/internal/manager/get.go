package manager

import "github.com/pkg/errors"

var ErrNodeNotFound = errors.New("node not found")

func (m *Manager) Get(id string) (*NodeCommand, error) {
	nodeCmd, ok := m.nodes[id]
	if !ok {
		return nil, ErrNodeNotFound
	}

	return nodeCmd, nil
}

func (m *Manager) GetAll() []*NodeCommand {
	nodes := make([]*NodeCommand, 0)

	for _, nodeCommand := range m.nodes {
		nodes = append(nodes, nodeCommand)
	}

	return nodes
}
