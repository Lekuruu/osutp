package database

import (
	"fmt"
	"time"
)

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
	ID                   int `gorm:"primaryKey"`
	SetID                int `gorm:"column:set_id;index"`
	Title                string
	Artist               string
	Creator              string
	Source               string
	Tags                 string
	Version              string               `gorm:"not null"`
	Status               int                  `gorm:"not null;default:1"`
	AR                   float64              `gorm:"not null"`
	OD                   float64              `gorm:"not null"`
	CS                   float64              `gorm:"not null"`
	AmountNormal         int                  `gorm:"not null"`
	AmountSliders        int                  `gorm:"not null"`
	AmountSpinners       int                  `gorm:"not null"`
	MaxCombo             int                  `gorm:"not null"`
	DifficultyAttributes DifficultyAttributes `gorm:"type:json"`
}

func (beatmap *Beatmap) FullName() string {
	return fmt.Sprintf("%s - %s (%s) [%s]", beatmap.Artist, beatmap.Title, beatmap.Creator, beatmap.Version)
}
