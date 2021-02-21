package models

import (
	"gorm.io/gorm"
)

// Rikishis - A rikishi array
type Rikishis []Rikishi

// Rikishi - Database model representing a Rikishi(Sumo Wrestler)
type Rikishi struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Avatar      string `gorm:"default:'/assets/default_avatar.jpg'"`
	East        bool
	SubRank     uint
	RankID      uint
	Rank        string `gorm:"-"`
	CategoryID  uint
	Category    string  `gorm:"-"`
	EastMatches []Match `gorm:"foreignKey:EastID"`
	WestMatches []Match `gorm:"foreignKey:WestID"`
}

// GetIDs - returns an array of the Rikishis' IDs
func (r Rikishis) GetIDs() []uint {
	rikishiIDs := make([]uint, len(r))
	for i, rikishi := range r {
		rikishiIDs[i] = rikishi.ID
	}
	return rikishiIDs
}
