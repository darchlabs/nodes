package main

import (
	"fmt"
	"log"

	"github.com/darchlabs/nodes/config"
	"github.com/darchlabs/nodes/internal/api"
	"github.com/darchlabs/nodes/internal/command"
	"github.com/darchlabs/nodes/internal/manager"
	"github.com/darchlabs/nodes/internal/storage"
	"github.com/darchlabs/nodes/pkg/namer"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	fmt.Println("------ Starting node runner")
	conf := &config.Config{}
	err := envconfig.Process("", conf)
	check(err)

	log.Printf("Database connection [darch node] done\n")
	store, err := storage.NewDataStore(conf.RedisURL)
	check(err)

	nameGenerator, err := namer.New()
	check(err)
	manager, err := manager.New(&manager.Config{
		MainConfig:    conf,
		IDGenerator:   uuid.NewString,
		NameGenerator: nameGenerator,
		// v1 config
		BootstrapNodesURL: conf.NetworksURL,
		BasePathDatabase:  conf.BasePathDatabase,
		// v2 config
		KubeConfigFilePath:  conf.KubeconfigFilePath,
		KubeconfigRemoteURL: conf.KubeconfigRemoteURL,
	})
	check(err)

	server := api.NewServer(&api.ServerConfig{
		Port:              conf.ApiServerPort,
		BootstrapNodesURL: conf.NetworksURL,
		Manager:           manager,
	})

	log.Printf("Starting [darch node]\n")
	err = server.Start(store)
	check(err)

	// listen interrupt
	quit := make(chan struct{})
	command.ListenInterruption(quit)
	<-quit
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
