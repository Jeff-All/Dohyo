package models

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var log *logrus.Logger
var db *gorm.DB

// GetModelFromID - Returns the model with the given string ID
func GetModelFromID(id string) (interface{}, error) {
	switch id {
	case "rank":
		return &Rank{}, nil
	case "rikishi":
		return &Rikishi{}, nil
	default:
		return nil, fmt.Errorf("Invalid model id '%s'", id)
	}
}
