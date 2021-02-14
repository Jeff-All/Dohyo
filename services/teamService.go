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

// NewTeamService - Instantiates a new TeamService
func NewTeamService(
	log *logrus.Logger,
	db *gorm.DB,
) TeamService {
	return TeamService{
		db:  db,
		log: log,
	}
}

// SaveRikishisToTeam - Saves the rikishi to the user's team
func (s *TeamService) SaveRikishisToTeam(user models.User, rikishis []models.Rikishi) error {
	team := models.Team{}
	var err error
	if err = s.db.Model(&user).Association("Team").Find(&team); err != nil {
		s.log.Errorf("error pulling team for user: %s", err)
		return err
	}
	if team.ID == 0 {
		s.log.Infof("team id is 0")
		if err = s.db.Create(&team).Error; err != nil {
			s.log.Errorf("error creating team: %s", err)
			return err
		}
		if err = s.db.Model(&user).Association("Team").Append(&team); err != nil {
			s.log.Errorf("error associating team to user: %s", err)
			return err
		}
	}
	if err = s.db.Model(&team).Association("Rikishis").Clear(); err != nil {
		s.log.Errorf("error while clearing team rikishi associations: %s", err)
		return err
	}
	if err = s.db.Model(&team).Association("Rikishis").Append(rikishis); err != nil {
		s.log.Errorf("error while appending rikishis to team")
		return err
	}
	return nil
}
