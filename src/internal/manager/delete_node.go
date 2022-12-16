package manager

func (m *Manager) Delete(id string) {
	delete(m.nodes, id)
}
