package manager

import (
	"fmt"

	"github.com/darchlabs/nodes/src/internal/command"
)

type IDGenerator func() string

type Manager struct {
	nodes                 map[string]*NodeCommand
	boostrapNodesURL      map[string]string
	idGenerator           IDGenerator
	currentAssignablePort int
	basePathDB            string
}

type Config struct {
	IDGenerator       IDGenerator
	BootstrapNodesURL map[string]string
	BasePathDatabase  string
}

func New(config *Config) *Manager {
	bootstrapNodesURL := make(map[string]string)
	for network, url := range config.BootstrapNodesURL {
		bootstrapNodesURL[network] = fmt.Sprintf("https://%s", url)
	}

	return &Manager{
		nodes:                 make(map[string]*NodeCommand),
		boostrapNodesURL:      bootstrapNodesURL,
		idGenerator:           config.IDGenerator,
		basePathDB:            config.BasePathDatabase,
		currentAssignablePort: 8545,
	}
}

type NodeCommand struct {
	ID     string
	Node   *command.Command
	Config *NodeConfig
}

type NodeConfig struct {
	Host              string
	Network           string
	Port              int
	BaseChainDataPath string
	BootsrapNodeURL   string
	FromBlockNumber   int64
	Label             string
}
