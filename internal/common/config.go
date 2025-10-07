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
	Server struct {
		Type   string `envconfig:"SERVER_TYPE" default:"titanic" validate:"oneof=titanic"`
		WebUrl string `envconfig:"SERVER_WEB_URL" default:"https://osu.titanic.sh"`
		ApiUrl string `envconfig:"SERVER_API_URL" default:"https://api.titanic.sh"`
	}
	TpWebsiteUrl string `envconfig:"TP_WEBSITE_URL" default:"https://tp.titanic.sh"`
}

func LoadConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
