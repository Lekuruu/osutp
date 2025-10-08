package common

import (
	"github.com/Lekuruu/osutp/internal/database"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Web struct {
		Host string `envconfig:"WEB_HOST" default:"0.0.0.0"`
		Port int    `envconfig:"WEB_PORT" default:"8080"`
	}
	Server struct {
		Type   string `envconfig:"SERVER_TYPE" default:"titanic" validate:"oneof=titanic"`
		WebUrl string `envconfig:"SERVER_WEB_URL" default:"https://osu.titanic.sh"`
		ApiUrl string `envconfig:"SERVER_API_URL" default:"https://api.titanic.sh"`
	}
	TpWebsiteUrl string `envconfig:"TP_WEBSITE_URL" default:"https://tp.titanic.sh"`
	Database     database.DatabaseConfig
}

func LoadConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
