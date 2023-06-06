package manager

func (m *Manager) DeleteInstance(id string) {
	delete(m.nodes, id)
}
