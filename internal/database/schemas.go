package database

import "time"

type Page struct {
	Name       string    `gorm:"primaryKey;not null"`
	Views      int64     `gorm:"not null;default:0"`
	LastUpdate time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type Changelog struct {
	Id          int       `gorm:"primaryKey;autoIncrement;not null"`
	Area        string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (changelog *Changelog) Date() string {
	return changelog.CreatedAt.UTC().Format("Jan 02, 2006")
}

func (changelog *Changelog) Time() string {
	return changelog.CreatedAt.UTC().Format("15:04")
}

type Beatmap struct {
	ID                   uint                 `gorm:"primaryKey"`
	SetID                int                  `gorm:"column:set_id;index"`
	Title                string               `gorm:"not null"`
	Artist               string               `gorm:"not null"`
	Creator              string               `gorm:"not null"`
	Source               string               `gorm:"not null"`
	Tags                 string               `gorm:"not null"`
	DifficultyAttributes DifficultyAttributes `gorm:"type:json"`
}
