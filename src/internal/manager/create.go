package manager

import (
	"fmt"
	"os"

	"github.com/darchlabs/nodes/src/internal/command"
	"github.com/pkg/errors"
)

type CreateNodeConfig struct {
	Network         string
	FromBlockNumber int64
}

func (m *Manager) CreateNode(config *CreateNodeConfig) (*NodeCommand, error) {
	id := m.idGenerator()
	bootstrapURL, ok := m.boostrapNodesURL[config.Network]
	if !ok {
		return nil, errors.New("manager: Manager.CreateNode network not supported")
	}

	nodeRunner, ok := networkNodesRunners[config.Network]
	if !ok {
		return nil, errors.New("manager: Manager.CreateNode network not supported")
	}

	dbPath := fmt.Sprintf("%s/%s/%d", m.basePathDB, config.Network, len(m.nodes))
	exist := existDir(dbPath)
	if !exist {
		mkdir := command.New("mkdir", "-p", dbPath)
		err := mkdir.Start()
		if err != nil {
			return nil, errors.Wrap(err, "manager: Manager.CreateNode mkdir.Start creating db dir error")
		}
	}

	nodeConfig := &NodeConfig{
		Network:           config.Network,
		Host:              "0.0.0.0",
		Port:              m.currentAssignablePort,
		BaseChainDataPath: dbPath,
		BootsrapNodeURL:   bootstrapURL,
		FromBlockNumber:   config.FromBlockNumber,
	}
	cmd := nodeRunner(nodeConfig)
	err := cmd.StreamOutput(id)
	if err != nil {
		return nil, errors.Wrap(err, "manager: Manager.CreateNode cmd.StreamOutput error")
	}

	err = cmd.Start()
	if err != nil {
		return nil, errors.Wrap(err, "manager: Manager.CreateNode cmd.Start error")
	}

	nodeCmd := &NodeCommand{
		ID:     id,
		Node:   cmd,
		Config: nodeConfig,
	}
	m.nodes[id] = nodeCmd

	m.currentAssignablePort++
	return nodeCmd, nil
}

func existDir(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return true
}
