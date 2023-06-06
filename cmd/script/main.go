package main

import (
	"fmt"

	"github.com/darchlabs/nodes/config"
	"github.com/darchlabs/nodes/internal/storage"
	_ "github.com/darchlabs/nodes/migrations"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	fmt.Println("------ Starting node scripts")
	conf := &config.Config{}
	err := envconfig.Process("", conf)
	check(err)
	// Initialize db connection
	storage, err := storage.NewSQLStore(conf.DBDriver, conf.PostgresDSN)
	check(err)

	err = goose.Up(storage.DB.DB, "migrations/")
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
