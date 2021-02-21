package services

import (
	"github.com/Jeff-All/Dohyo/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// MatchService - Service for handling matches
type MatchService struct {
	db  *gorm.DB
	log *logrus.Logger

	rikishiService    RikishiService
	tournamentService TournamentService
}

// NewMatchService - Instantiates a new MatchService
func NewMatchService(
	log *logrus.Logger,
	db *gorm.DB,
	rikishiService RikishiService,
	tournamentService TournamentService,
) MatchService {
	return MatchService{
		log:               log,
		db:                db,
		rikishiService:    rikishiService,
		tournamentService: tournamentService,
	}
}

// AddMatches - Adds the provided matches to the database
func (s *MatchService) AddMatches(matches []models.Match) error {
	s.log.Infof("adding %d matches", len(matches))

	var err error
	var rikishis map[string]models.Rikishi
	if rikishis, err = s.rikishiService.GetRikishiMappedByName(); err != nil {
		s.log.Errorf("error getting rikishis by name: %s", err)
		return err
	}
	if err = s.db.Create(matches).Error; err != nil {
		s.log.Errorf("error creating matches: %s", err)
		return err
	}
	var tournaments map[string]models.Tournament
	if tournaments, err = s.tournamentService.GetTournamentsByName(); err != nil {
		s.log.Errorf("error getting tournaments by name: %s", err)
		return err
	}
	for _, match := range matches {
		if tournament, ok := tournaments[match.Tournament]; !ok {
			s.log.Infof("unknown tournament '%s'", match.Tournament)
		} else {
			s.log.Infof("binding tournament %s to match", match.Tournament)
			s.db.Model(&tournament).Association("Matches").Append(&match)
		}
		if rikishi, ok := rikishis[match.East]; !ok {
			s.log.Infof("unknown rikishi '%s' for match", match.East)
		} else {
			s.log.Infof("binding rikishi '%s' to match", match.East)
			if err = s.db.Model(&rikishi).Association("EastMatches").Append(&match); err != nil {
				s.log.Errorf("error while binding rikishi '%s' to match: %s", match.East, err)
				return err
			}
		}
		if rikishi, ok := rikishis[match.West]; !ok {
			s.log.Infof("unknown rikishi '%s' for match", match.West)
		} else {
			s.log.Infof("binding rikishi '%s' to match", match.West)
			if err = s.db.Model(&rikishi).Association("WestMatches").Append(&match); err != nil {
				s.log.Errorf("error while binding rikishi '%s' to match: %s", match.West, err)
				return err
			}
		}
	}
	for _, tournament := range tournaments {
		if err = s.db.Save(&tournament).Error; err != nil {
			s.log.Errorf("error while saving tournament updates: %s", err)
			return err
		}
	}
	for _, rikishi := range rikishis {
		if err = s.db.Save(&rikishi).Error; err != nil {
			s.log.Errorf("error while saving rikishi updates: %s", err)
			return err
		}
	}

	s.log.Infof("successfully added %d matches", len(matches))
	return nil
}
