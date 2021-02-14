package models

import "gorm.io/gorm"

// User - User model
type User struct {
	gorm.Model
	Auth0ID string
	Email   string
	Team    Team
}
