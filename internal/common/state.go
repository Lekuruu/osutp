package common

import (
	"log"
	"os"

	"github.com/Lekuruu/osutp/internal/database"
	"gorm.io/gorm"
)

type State struct {
	Config   *Config
	Database *gorm.DB
	Logger   *Logger
}

func NewState() *State {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return nil
	}

	// Ensure data path exists
	if err := os.MkdirAll(".data", os.ModePerm); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
		return nil
	}

	db, err := database.CreateSession(config.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
		return nil
	}

	return &State{
		Logger:   NewLogger("osutp"),
		Config:   config,
		Database: db,
	}
}
