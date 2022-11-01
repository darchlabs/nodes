package main

import (
	"fmt"
	"log"

	"github.com/darchlabs/nodes/src/config"
	"github.com/darchlabs/nodes/src/internal/api"
	"github.com/darchlabs/nodes/src/internal/command"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	fmt.Println("------ Starting node runner")
	var conf config.Config
	err := envconfig.Process("", &conf)
	check(err)

	log.Printf("Starting [darch %s node]\n", conf.Chain)
	log.Printf("%s", conf.BaseChainDataPath)

	nodeURL := fmt.Sprintf("%s@%s", conf.NodeURL, conf.BlockNumber)

	server := api.NewServer(&api.ServerConfig{
		Port: conf.ApiServerPort,
		CommandConfig: &api.CommandConfig{
			Runner:           "ganache",
			Host:             "0.0.0.0",
			DatabasePath:     fmt.Sprintf("%s", conf.BaseChainDataPath),
			BootstrapNodeURL: nodeURL,
		},
	})

	err = server.Start()
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
