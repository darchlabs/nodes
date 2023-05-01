package manager

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/darchlabs/nodes/config"
	"github.com/darchlabs/nodes/internal/command"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var ErrKeyNotFound = errors.New("environment key not found")

type IDGenerator func() string

type NameGenerator interface {
	Generate() string
}

type Manager struct {
	MainConfig            *config.Config
	nodes                 map[string]*NodeInstance
	boostrapNodesURL      map[string]string
	idGenerator           IDGenerator
	nameGenerator         NameGenerator
	currentAssignablePort int
	basePathDB            string
	networkNodeSetups     map[string]nodeSetup
	// v2 related
	clusterClient *kubernetes.Clientset
}

type Config struct {
	MainConfig    *config.Config
	IDGenerator   IDGenerator
	NameGenerator NameGenerator
	// v1 config related
	BootstrapNodesURL map[string]string
	BasePathDatabase  string
	// v2 config related
	KubeConfigFilePath  string
	KubeconfigRemoteURL string
}

func New(config *Config) (*Manager, error) {
	bootstrapNodesURL := make(map[string]string)
	for network, url := range config.BootstrapNodesURL {
		bootstrapNodesURL[network] = fmt.Sprintf("https://%s", url)
	}

	// v2 kubernetes setup
	//get remote file if exist. Otherwise path only will be used
	if config.KubeconfigRemoteURL != "" {
		fmt.Println("--------- url file ", config.KubeconfigRemoteURL)
		res, err := http.Get(config.KubeconfigRemoteURL)
		if err != nil {
			return nil, errors.Wrap(err, "manager: New http.Get error")
		}

		out, err := os.Create(config.KubeConfigFilePath)
		if err != nil {
			return nil, errors.Wrap(err, "manager: New http.Get ")
		}

		// Copy the contents of the response body to the output file
		_, err = io.Copy(out, res.Body)
		if err != nil {
			return nil, errors.Wrap(err, "manager: New io.Copy error")
		}

		if err = out.Close(); err != nil {
			return nil, errors.Wrap(err, "manager: New body.Close error")
		}
	}
	// using the file created
	k8sClusterConfig, err := clientcmd.BuildConfigFromFlags("", config.KubeConfigFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "manager: New clientcmd.BuildConfigFromFlags error")
	}

	clusterClient, err := kubernetes.NewForConfig(k8sClusterConfig)
	if err != nil {
		return nil, errors.Wrap(err, "manager: New kubernetes.NewForConfig error")
	}

	config.MainConfig.Images = config.MainConfig.ParseImages()
	for k, v := range config.MainConfig.Images {
		fmt.Println("images for", k, "related images", v)
	}
	m := &Manager{
		MainConfig:            config.MainConfig,
		nodes:                 make(map[string]*NodeInstance),
		boostrapNodesURL:      bootstrapNodesURL,
		idGenerator:           config.IDGenerator,
		nameGenerator:         config.NameGenerator,
		basePathDB:            config.BasePathDatabase,
		currentAssignablePort: 8545,
		clusterClient:         clusterClient,
	}
	m.networkNodeSetups = setupFuncByNetwork(m)

	return m, nil
}

type NodeInstance struct {
	ID        string
	Name      string
	Node      *command.Command
	Config    *NodeConfig
	Artifacts *Artifacts
}

type NodeConfig struct {
	Host              string
	Network           string
	Environment       string
	Port              int
	BaseChainDataPath string
	BootsrapNodeURL   string
	FromBlockNumber   int64
	Label             string
	CreatedAt         time.Time
}

type Artifacts struct {
	Deployments []string
	Pods        []string
	Services    []string
}

type nodeSetup func(network string, env map[string]string) (*NodeInstance, error)
