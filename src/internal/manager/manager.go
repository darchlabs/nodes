package manager

import (
	"fmt"
	"time"

	"github.com/darchlabs/nodes/src/internal/command"
)

type IDGenerator func() string

type NameGenerator interface {
	Generate() string
}

type Manager struct {
	nodes                 map[string]*NodeInstance
	boostrapNodesURL      map[string]string
	idGenerator           IDGenerator
	nameGenerator         NameGenerator
	currentAssignablePort int
	basePathDB            string
}

type Config struct {
	IDGenerator       IDGenerator
	NameGenerator     NameGenerator
	BootstrapNodesURL map[string]string
	BasePathDatabase  string
}

func New(config *Config) *Manager {
	bootstrapNodesURL := make(map[string]string)
	for network, url := range config.BootstrapNodesURL {
		bootstrapNodesURL[network] = fmt.Sprintf("https://%s", url)
	}

	return &Manager{
		nodes:                 make(map[string]*NodeInstance),
		boostrapNodesURL:      bootstrapNodesURL,
		idGenerator:           config.IDGenerator,
		nameGenerator:         config.NameGenerator,
		basePathDB:            config.BasePathDatabase,
		currentAssignablePort: 8545,
	}
}

type NodeInstance struct {
	ID     string
	Name   string
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
	CreatedAt         time.Time
}
