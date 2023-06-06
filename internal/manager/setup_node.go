package manager

func setupFuncByNetwork(m *Manager) map[string]nodeSetup {
	return map[string]nodeSetup{
		"ethereum":  m.EvmDevNode,
		"polygon":   m.EvmDevNode,
		"chainlink": m.ChainlinkNode,
		"celo":      m.CeloNode,
	}
}
