package models

import "gorm.io/gorm"

// Tournament - Represents a tournament
type Tournament struct {
	gorm.Model
	Name  string `gorm:"unique"`
	Month string
	Year  uint
}
