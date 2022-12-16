package main

import (
	"fmt"
	"log"

	"github.com/darchlabs/nodes/src/config"
	"github.com/darchlabs/nodes/src/internal/api"
	"github.com/darchlabs/nodes/src/internal/command"
	"github.com/darchlabs/nodes/src/internal/manager"
	"github.com/darchlabs/nodes/src/internal/storage"
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

	log.Printf("Starting [darch node]\n")

	server := api.NewServer(&api.ServerConfig{
		Port:              conf.ApiServerPort,
		BootstrapNodesURL: conf.NetworksURL,
		Manager: manager.New(&manager.Config{
			IDGenerator:       uuid.NewString,
			BootstrapNodesURL: conf.NetworksURL,
			BasePathDatabase:  conf.BasePathDatabase,
		}),
	})

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
