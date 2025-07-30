package database

import "time"

type Page struct {
	Name       string    `gorm:"primaryKey;not null"`
	Views      int64     `gorm:"not null;default:0"`
	LastUpdate time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
