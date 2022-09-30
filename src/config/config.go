package config

type Config struct {
	Environment string `envconfig:"environment" required:"true"`
	Chain       string `envconfig:"chain" required:"true"`
	NodeURL     string `envconfig:"node_url" required:"true"`
}
