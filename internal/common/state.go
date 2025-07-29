package common

import (
	"log"

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
		log.Fatalf("Failed to load configuration: %v", err)
		return nil
	}

	db, err := database.CreateSession(config.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
		return nil
	}

	return &State{
		Config:   config,
		Database: db,
	}
}
