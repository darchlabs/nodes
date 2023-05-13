package application

import (
	"github.com/darchlabs/nodes/config"
	"github.com/darchlabs/nodes/internal/manager"
	"github.com/darchlabs/nodes/internal/storage"
	"github.com/darchlabs/nodes/pkg/namer"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type App struct {
	Manager       *manager.Manager
	KeyValueStore storage.KeyValue
	SqlStore      storage.SQL
	NameGenerator *namer.Namer
}

type Config struct {
	KeyValueStore storage.KeyValue
	SqlStore      storage.SQL
	MainConfig    *config.Config
}

func NewApp(conf *Config) (*App, error) {
	nameGen, err := namer.New()
	if err != nil {
		return nil, errors.Wrap(err, "app: NewApp namer.New error")
	}

	manager, err := manager.New(&manager.Config{
		MainConfig:    conf.MainConfig,
		IDGenerator:   uuid.NewString,
		NameGenerator: nameGen,
		// v1 config
		BootstrapNodesURL: conf.MainConfig.NetworksURL,
		BasePathDatabase:  conf.MainConfig.BasePathDatabase,
		// v2 config
		KubeConfigFilePath:  conf.MainConfig.KubeconfigFilePath,
		KubeconfigRemoteURL: conf.MainConfig.KubeconfigRemoteURL,
	})
	if err != nil {
		return nil, errors.Wrap(err, "app: NewApp manager.New error")
	}

	return &App{
		Manager:       manager,
		KeyValueStore: conf.KeyValueStore,
		SqlStore:      conf.SqlStore,
		NameGenerator: nameGen,
	}, nil
}
