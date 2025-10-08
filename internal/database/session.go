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
			SlowThreshold:             500 * time.Millisecond,
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

	err = CreateIndexes(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateIndexes(db *gorm.DB) error {
	indexStatements := []string{
		`CREATE INDEX IF NOT EXISTS idx_beatmaps_star_rating
		 ON beatmaps (json_extract(difficulty_attributes, '$.0.StarRating'));`,

		`CREATE INDEX IF NOT EXISTS idx_beatmaps_speed_stars
		 ON beatmaps (json_extract(difficulty_attributes, '$.0.SpeedStars'));`,

		`CREATE INDEX IF NOT EXISTS idx_beatmaps_aim_stars
		 ON beatmaps (json_extract(difficulty_attributes, '$.0.AimStars'));`,
	}

	for _, stmt := range indexStatements {
		if err := db.Exec(stmt).Error; err != nil {
			return err
		}
	}

	return nil
}

func PreloadQuery(database *gorm.DB, preload []string) *gorm.DB {
	result := database

	for _, p := range preload {
		result = result.Preload(p)
	}

	return result
}
