package common

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Web struct {
		Host string `envconfig:"API_HOST" default:"0.0.0.0"`
		Port int    `envconfig:"API_PORT" default:"8080"`
	}
	Database struct {
		Path string `envconfig:"DB_PATH" default:"./.data/osutp.db"`
	}
}

func LoadConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
