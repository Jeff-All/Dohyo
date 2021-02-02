package models

// Rank - Rikishi ranks
type Rank struct {
	ID       uint `gorm:"primaryKey"`
	Level    uint
	Name     string
	Rikishis []Rikishi `gorm:"foreignKey:Rank"`
}
