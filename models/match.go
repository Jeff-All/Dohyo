package models

import "gorm.io/gorm"

// Match - Model representing a single sumo match
type Match struct {
	gorm.Model
	TournamentID uint
	Day          uint
	EastID       uint
	WestID       uint
	WinnerID     uint
	Winner       string `gorm:"-"`
	Tournament   string `gorm:"-"`
	East         string `gorm:"-"`
	West         string `gorm:"-"`
}
