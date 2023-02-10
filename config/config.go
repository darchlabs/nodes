package config

type Config struct {
	Environment string `envconfig:"environment" required:"true"`

	// node config
	BlockNumber      string            `envconfig:"block_number" default:"1"`
	BasePathDatabase string            `envconfig:"base_path_database" default:"/data"`
	NetworksURL      map[string]string `envconfig:"networks_url" required:"true"`

	// server config
	ApiServerHost string `envconfig:"api_server_host" default:"0.0.0.0."`
	ApiServerPort string `envconfig:"api_server_port" default:"6969"`
	MasterURL     string `envconfig:"-" default:"http://master.darchlabs.com/nodes/status"`

	// database config
	RedisURL string `envconfig:"redis_url" required:"true"`
}
