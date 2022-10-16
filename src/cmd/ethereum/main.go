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
	var conf config.Config
	err := envconfig.Process("", &conf)
	check(err)

	log.Printf("Starting [darch %s node]\n", conf.Chain)

	cmd := command.New(
		"ganache",
		"--host", "0.0.0.0",
		"--db", "/home/node/data",
		"--fork", fmt.Sprintf("%s@%s", conf.NodeURL, conf.BlockNumber),
	)

	log.Println("Running command : ", cmd.Slug())

	server := api.NewServer(&api.ServerConfig{
		Port:    conf.ApiServerPort,
		Command: cmd,
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
