package manager

import (
	"fmt"
	"time"

	"github.com/darchlabs/nodes/internal/command"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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
	// v2 related
	clusterClient *kubernetes.Clientset
}

type Config struct {
	IDGenerator   IDGenerator
	NameGenerator NameGenerator
	// v1 config related
	BootstrapNodesURL map[string]string
	BasePathDatabase  string
	// v2 config related
	KubeConfigFilePath string
}

func New(config *Config) (*Manager, error) {
	bootstrapNodesURL := make(map[string]string)
	for network, url := range config.BootstrapNodesURL {
		bootstrapNodesURL[network] = fmt.Sprintf("https://%s", url)
	}

	// v2 kubernetes setup
	k8sClusterConfig, err := clientcmd.BuildConfigFromFlags("", config.KubeConfigFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "manager: New clientcmd.BuildConfigFromFlags error")
	}

	clusterClient, err := kubernetes.NewForConfig(k8sClusterConfig)
	if err != nil {
		return nil, errors.Wrap(err, "manager: New kubernetes.NewForConfig error")
	}

	return &Manager{
		nodes:                 make(map[string]*NodeInstance),
		boostrapNodesURL:      bootstrapNodesURL,
		idGenerator:           config.IDGenerator,
		nameGenerator:         config.NameGenerator,
		basePathDB:            config.BasePathDatabase,
		currentAssignablePort: 8545,
		clusterClient:         clusterClient,
	}, nil
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
