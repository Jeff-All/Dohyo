package models

import (
	"gorm.io/gorm"
)

// Rikishi - Database model representing a Rikishi(Sumo Wrestler)
type Rikishi struct {
	gorm.Model
	Name   string `gorm:"unique"`
	Avatar string `gorm:"default:'/assets/default_avatar.jpg'"`
	RankID uint
	Rank   string `gorm:"-"`
}
