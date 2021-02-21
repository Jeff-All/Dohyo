package models

import "gorm.io/gorm"

type Tournaments []Tournament

// Tournament - Represents a tournament
type Tournament struct {
	gorm.Model
	Name    string `gorm:"unique"`
	Month   string
	Year    uint
	Current bool
	Matches []Match
}

// MapByName - Returns a map of tournaments indexed by the 'name' column
func (t Tournaments) MapByName() map[string]Tournament {
	toReturn := make(map[string]Tournament, len(t))
	for _, tournament := range t {
		toReturn[tournament.Name] = tournament
	}
	return toReturn
}
