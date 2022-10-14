package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/darchlabs/nodes/src/config"
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
	err = cmd.Start()
	check(err)

	log.Println(cmd.Status())

	// listen interrupt
	quit := make(chan struct{})
	listenInterrupt(quit)
	<-quit
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func listenInterrupt(quit chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-c
		fmt.Println("Signal received", s.String())
		quit <- struct{}{}
	}()
}
