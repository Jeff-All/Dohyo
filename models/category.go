package models

// Category - Contest Categories
type Category struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	Rikishis     []Rikishi `gorm:"foreignKey:CategoryID"`
	RikishiNames []string  `gorm:"-"`
}
