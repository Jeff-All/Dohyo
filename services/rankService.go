package services

import (
	"errors"

	"github.com/Jeff-All/Dohyo/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RankService - Provides methods for handling Rikishi Ranks
type RankService struct {
	log *logrus.Logger
	db  *gorm.DB
}

// NewRankService - Instantiates a new RankService
func NewRankService(
	log *logrus.Logger,
	db *gorm.DB,
) RankService {
	return RankService{
		log: log,
		db:  db,
	}
}

// GetRanks - Returns all ranks from the database
func (s *RankService) GetRanks() ([]models.Rank, error) {
	ranks := []models.Rank{}
	var result *gorm.DB
	if result = s.db.Find(&ranks); result.Error == nil {
		return ranks, nil
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		s.log.Infof("unable to find any rank entries")
		return ranks, nil
	}
	s.log.Errorf("error while retrieving rank entries: %s", result.Error)
	return ranks, result.Error
}

// GetRanksMappedByName - Returns all ranks from the DB mapped by the name column
func (s *RankService) GetRanksMappedByName() (map[string]models.Rank, error) {
	var err error
	var ranks []models.Rank
	if ranks, err = s.GetRanks(); err != nil {
		return nil, err
	}
	rankMap := make(map[string]models.Rank)
	for _, rank := range ranks {
		rankMap[rank.Name] = rank
	}
	return rankMap, nil
}
