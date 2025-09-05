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
	CreatedAt            time.Time            `gorm:"not null;default:CURRENT_TIMESTAMP"`
	LastScoreUpdate      time.Time            `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DifficultyAttributes DifficultyAttributes `gorm:"type:json;default:null"`
}

func (beatmap *Beatmap) FullName() string {
	return fmt.Sprintf("%s - %s (%s) [%s]", beatmap.Artist, beatmap.Title, beatmap.Creator, beatmap.Version)
}

func (beatmap *Beatmap) ApproachRate(mods uint32) float64 {
	return beatmap.DifficultyAttributes[mods]["ApproachRate"]
}

func (beatmap *Beatmap) OverallDifficulty(mods uint32) float64 {
	return beatmap.DifficultyAttributes[mods]["OverallDifficulty"]
}

func (beatmap *Beatmap) CircleSize(mods uint32) float64 {
	return beatmap.DifficultyAttributes[mods]["CircleSize"]
}

func (beatmap *Beatmap) StarRating(mods uint32) float64 {
	return beatmap.DifficultyAttributes[mods]["StarRating"]
}

func (beatmap *Beatmap) SpeedStars(mods uint32) float64 {
	return beatmap.DifficultyAttributes[mods]["SpeedStars"]
}

func (beatmap *Beatmap) AimStars(mods uint32) float64 {
	return beatmap.DifficultyAttributes[mods]["AimStars"]
}

type Player struct {
	ID               int       `gorm:"primaryKey;autoIncrement;not null"`
	Name             string    `gorm:"not null;uniqueIndex"`
	Country          string    `gorm:"not null;default:'XX'"`
	GlobalRank       int       `gorm:"not null;default:0;index"`
	CountryRank      int       `gorm:"not null;default:0;index"`
	RecentRankChange int       `gorm:"not null;default:0"`
	TotalTp          float64   `gorm:"not null;default:0;index"`
	AimTp            float64   `gorm:"not null;default:0"`
	SpeedTp          float64   `gorm:"not null;default:0"`
	AccuracyTp       float64   `gorm:"not null;default:0"`
	CreatedAt        time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	LastUpdate       time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type Score struct {
	ID         int       `gorm:"primaryKey;autoIncrement;not null"`
	BeatmapID  int       `gorm:"not null;index"`
	PlayerID   int       `gorm:"not null;index"`
	Checksum   string    `gorm:"not null;size:32"`
	TotalScore int64     `gorm:"not null"`
	MaxCombo   int       `gorm:"not null"`
	Mods       uint32    `gorm:"not null;default:0"`
	FullCombo  bool      `gorm:"not null;default:false"`
	Grade      string    `gorm:"not null;default:'N';size:2"`
	Accuracy   float64   `gorm:"not null"`
	Amount300  int       `gorm:"not null"`
	Amount100  int       `gorm:"not null"`
	Amount50   int       `gorm:"not null"`
	AmountGeki int       `gorm:"not null"`
	AmountKatu int       `gorm:"not null"`
	AmountMiss int       `gorm:"not null"`
	TotalTp    float64   `gorm:"not null;default:0;index"`
	AimTp      float64   `gorm:"not null;default:0;index"`
	SpeedTp    float64   `gorm:"not null;default:0;index"`
	AccTp      float64   `gorm:"not null;default:0;index"`
	CreatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	LastUpdate time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
