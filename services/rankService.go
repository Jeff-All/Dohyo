package services

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RankService - Provides methods for handling Rikishi Ranks
type RankService struct {
	log logrus.Logger
	db  *gorm.DB
}
