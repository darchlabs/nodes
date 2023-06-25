package config

import (
	"os"
	"strings"
)

type Config struct {
	Environment string `envconfig:"environment" required:"true"`

	// node config
	BlockNumber      string            `envconfig:"block_number" default:"1"`
	BasePathDatabase string            `envconfig:"base_path_database" default:"/data"`
	NetworksURL      map[string]string `envconfig:"networks_url" required:"false"`

	// server config
	ApiServerHost    string `envconfig:"api_server_host" default:"0.0.0.0."`
	ApiServerPort    string `envconfig:"api_server_port" default:"6969"`
	MasterURL        string `envconfig:"-" default:"http://master.darchlabs.com/nodes/status"`
	BackofficeApiURL string `envconfig:"backoffice_api_url" required:"true"`

	// database config
	RedisURL    string `envconfig:"redis_url" required:"true"`
	DBDriver    string `envconfig:"db_driver" default:"postgres"`
	PostgresDSN string `envconfig:"postgres_dsn" required:"true"`

	// kubernetes config
	KubeconfigFilePath  string `envconfig:"kubeconfig_file_path" required:"true"`
	KubeconfigRemoteURL string `envconfig:"kubeconfig_remote_url" required:"false"`

	// images supported
	Images map[string]string
}

func (c *Config) ParseImages() map[string]string {
	imgs := os.Getenv("IMAGES_SUPPORTED")

	images := make(map[string]string)

	for _, pair := range strings.Split(imgs, ",") {
		// network $ image:version
		values := strings.Split(pair, "|")
		images[values[0]] = values[1]
	}
	return images
}
