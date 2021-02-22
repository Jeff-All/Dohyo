package services

import (
	"github.com/Jeff-All/Dohyo/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// TournamentService - Service for handling tournaments
type TournamentService struct {
	log *logrus.Logger
	db  *gorm.DB

	currentTournament *models.Tournament
}

// NewTournamentService - Instantiates a new TournamentService
func NewTournamentService(
	log *logrus.Logger,
	db *gorm.DB,
) TournamentService {
	return TournamentService{
		log:               log,
		db:                db,
		currentTournament: &models.Tournament{},
	}
}

// SetCurrentTournament - Sets the current active tournament in the database
func (s *TournamentService) SetCurrentTournament(name string) error {
	s.log.Infof("setting current tournament to '%s'", name)
	if err := s.db.Model(models.Tournament{}).Where("1=1").UpdateColumn("current", false).Error; err != nil {
		s.log.Errorf("error while trying to deactivate tournaments: %s", err)
		return err
	}
	tournament := models.Tournament{}
	if err := s.db.Where("name = ?", name).First(&tournament).Error; err != nil {
		s.log.Errorf("error while pulling tournament '%s': %s", name, err)
		return err
	}
	tournament.Current = true
	if err := s.db.Save(tournament).Error; err != nil {
		s.log.Errorf("error while saving current tournament: %s", err)
		return err
	}
	s.currentTournament = &tournament
	return nil
}

// GetCurrentTournament - Retrieves the current tournament from the TournamentService.currentTournament
// if null it pulls and sets it from the database
func (s *TournamentService) GetCurrentTournament() (*models.Tournament, error) {
	if s.currentTournament.ID == 0 {
		if err := s.db.Where("current = true").First(s.currentTournament).Error; err != nil {
			s.log.Errorf("error while retrieving current tournament: %s", err)
			return nil, err
		}
		s.log.Infof("pulled current tournament '%d'", s.currentTournament.ID)
	}
	return s.currentTournament, nil
}

// GetAllTournaments - Retrieves all tournaments in the database
func (s *TournamentService) GetAllTournaments() (models.Tournaments, error) {
	tournaments := models.Tournaments{}
	var err error
	if err = s.db.Find(&tournaments).Error; err != nil {
		s.log.Errorf("error pulling all tournaments: %s", err)
		return nil, err
	}
	return tournaments, nil
}

// GetTournamentsByName - Retrieves a map of tournaments mapped by their names
func (s *TournamentService) GetTournamentsByName() (map[string]models.Tournament, error) {
	var tournaments models.Tournaments
	var err error
	if tournaments, err = s.GetAllTournaments(); err != nil {
		s.log.Errorf("error while getting all tournaments: %s", err)
		return nil, err
	}
	return tournaments.MapByName(), nil
}
