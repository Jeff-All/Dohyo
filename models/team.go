package models

import "gorm.io/gorm"

// Team - Model of a user's team
type Team struct {
	gorm.Model
	UserID   uint
	Rikishis []Rikishi `gorm:"many2many:team_rikishis;"`
}

// IsCreated - returns whether a team has been instantiated in the DB
func (t Team) IsCreated() bool {
	return t.ID != 0
}
