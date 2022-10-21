package main

import (
	"fmt"
	"log"

	"github.com/darchlabs/nodes/src/config"
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

	cmd := command.New(
		"ganache",
		"--host", "0.0.0.0",
		"--db", fmt.Sprintf("%s", conf.BaseChainDataPath),
		"--fork", nodeURL,
	)

	//server := api.NewServer(&api.ServerConfig{
	//Port: conf.ApiServerPort,
	//Command: command.New(
	//"ganache",
	//"--host", "0.0.0.0",
	//"--db", fmt.Sprintf("%s", conf.BaseChainDataPath),
	//"--fork", nodeURL,
	//)})

	//err = server.Start()
	err = cmd.Start()
	check(err)

	//log.Println(cmd.Status())

	check(err)
	//log.Println("Running command : ", cmd.Slug())
	//log.Println("Running command : ", cmd.Status())

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
