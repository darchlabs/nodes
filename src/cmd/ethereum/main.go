package main

import (
	"fmt"
	"log"

	"github.com/darchlabs/nodes/src/config"
	"github.com/darchlabs/nodes/src/internal/api"
	"github.com/darchlabs/nodes/src/internal/command"
	"github.com/darchlabs/nodes/src/internal/manager"
	"github.com/darchlabs/nodes/src/internal/storage"
	"github.com/darchlabs/nodes/src/pkg/namer"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	fmt.Println("------ Starting node runner")
	var conf config.Config
	err := envconfig.Process("", &conf)
	check(err)

	log.Printf("Database connection [darch node] done\n")
	store, err := storage.NewDataStore(conf.RedisURL)
	check(err)

	nameGenerator, err := namer.New()
	check(err)

	server := api.NewServer(&api.ServerConfig{
		Port:              conf.ApiServerPort,
		BootstrapNodesURL: conf.NetworksURL,
		Manager: manager.New(&manager.Config{
			IDGenerator:       uuid.NewString,
			NameGenerator:     nameGenerator,
			BootstrapNodesURL: conf.NetworksURL,
			BasePathDatabase:  conf.BasePathDatabase,
		}),
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
