package services

import (
	"github.com/Jeff-All/Dohyo/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// TeamService - Service for handling teams
type TeamService struct {
	db  *gorm.DB
	log *logrus.Logger
}

// SaveTeam - Saves the team
func (s *TeamService) SaveTeam(user *models.User, team *models.Team) {

}
