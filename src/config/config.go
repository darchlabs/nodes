package config

type Config struct {
	Environment string `envconfig:"environment" required:"true"`

	// node config
	Chain             string `envconfig:"chain" required:"true"`
	NodeURL           string `envconfig:"node_url" required:"true"`
	BlockNumber       string `envconfig:"block_number" default:"1"`
	BaseChainDataPath string `envconfig:"base_chain_data_path" default:"/data"`

	// server config
	ApiServerHost string `envconfig:"api_server_host" default:"0.0.0.0."`
	ApiServerPort string `envconfig:"api_server_port" default:"6969"`
	MasterURL     string `envconfig:"-" default:"http://master.darchlabs.com/nodes/status"`
}
