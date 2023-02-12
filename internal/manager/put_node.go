package manager

func (m *Manager) Put(nc *NodeInstance) {
	m.nodes[nc.ID] = nc
}
