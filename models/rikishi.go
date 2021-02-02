package models

import (
	"gorm.io/gorm"
)

// Rikishi - Database model representing a Rikishi(Sumo Wrestler)
type Rikishi struct {
	gorm.Model
	Name string `gorm:"unique"`
	Rank uint
}
