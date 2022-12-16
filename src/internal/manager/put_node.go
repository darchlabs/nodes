package manager

func (m *Manager) Put(nc *NodeCommand) {
	m.nodes[nc.ID] = nc
}
