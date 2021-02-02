package services

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RikishiService - Provides functions for accessing the Rikishi data
type RikishiService struct {
	log logrus.Logger
	db  *gorm.DB
}
