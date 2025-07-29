package database

type Page struct {
	Name  string `gorm:"primaryKey;not null"`
	Views int64  `gorm:"not null;default:0"`
}
