package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateSession(path string) (*gorm.DB, error) {
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	var models = []interface{}{
		&Page{},
		&Changelog{},
		&Beatmap{},
		&Player{},
		&Score{},
	}

	err = db.AutoMigrate(models...)
	if err != nil {
		return nil, err
	}
	return db, nil
}
