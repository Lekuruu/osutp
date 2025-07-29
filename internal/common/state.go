package common

import (
	"github.com/Lekuruu/osutp-web/internal/database"
	"gorm.io/gorm"
)

type State struct {
	Config   *Config
	Database *gorm.DB
}

func NewState() *State {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := database.CreateSession(config.Database.Path)
	if err != nil {
		panic(err)
	}

	return &State{
		Config:   config,
		Database: db,
	}
}
