package main

import (
	"fmt"
	"log"

	"github.com/darchlabs/nodes/config"
	"github.com/darchlabs/nodes/internal/api"
	"github.com/darchlabs/nodes/internal/application"
	"github.com/darchlabs/nodes/internal/command"
	"github.com/darchlabs/nodes/internal/storage"
	_ "github.com/darchlabs/nodes/migrations"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	fmt.Println("------ Starting node runner")
	conf := &config.Config{}
	err := envconfig.Process("", conf)
	check(err)

	log.Printf("Database postgresql connection [darch node] done\n")
	sqlStore, err := storage.NewSQLStore(conf.DBDriver, conf.PostgresDSN)
	check(err)

	err = goose.Up(sqlStore.DB.DB, "migrations/")
	check(err)

	log.Printf("Database Redis connection [darch node] done\n")
	kvStore, err := storage.NewKeyValueStore(conf.RedisURL)
	check(err)

	app, err := application.NewApp(&application.Config{
		SqlStore:      sqlStore,
		KeyValueStore: kvStore,
		MainConfig:    conf,
	})
	check(err)
	fmt.Println(app)

	server := api.NewServer(&api.ServerConfig{
		Port:              conf.ApiServerPort,
		BootstrapNodesURL: conf.NetworksURL,
		Manager:           app.Manager,
		App:               app,
	})

	log.Printf("Starting [darch node]\n")
	err = server.Start(app)
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
