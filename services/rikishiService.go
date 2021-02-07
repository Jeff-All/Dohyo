package services

import (
	"errors"

	"github.com/Jeff-All/Dohyo/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RikishiService - Provides functions for accessing the Rikishi data
type RikishiService struct {
	log         *logrus.Logger
	db          *gorm.DB
	rankService RankService
}

// NewRikishiService - Instantiates a new RikishiService
func NewRikishiService(
	log *logrus.Logger,
	db *gorm.DB,
	rankService RankService,
) RikishiService {
	return RikishiService{
		log:         log,
		db:          db,
		rankService: rankService,
	}
}

// GetAllRikishi - Returns all rikishi in the database
func (s *RikishiService) GetAllRikishi() ([]models.Rikishi, error) {
	rikishis := []models.Rikishi{}
	var result *gorm.DB
	if result = s.db.Find(&rikishis); result.Error == nil {
		s.log.Infof("successfully pulled all rikishi")
		return rikishis, nil
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		s.log.Infof("unable to find any rikishi entries")
		return rikishis, nil
	}
	s.log.Errorf("error while retrieving rikishi entries: %s", result.Error)
	return rikishis, result.Error
}

// GetRikishiMappedByName - Returns a map of all rikishi indexed by their Name column
func (s *RikishiService) GetRikishiMappedByName() (map[string]models.Rikishi, error) {
	var rikishis []models.Rikishi
	var err error
	if rikishis, err = s.GetAllRikishi(); err != nil {
		s.log.Errorf("error while pulling rikishi: %s", err)
		return nil, err
	}
	rikishiMap := make(map[string]models.Rikishi)
	for _, rikishi := range rikishis {
		rikishiMap[rikishi.Name] = rikishi
	}
	return rikishiMap, nil
}

// AddRikishi - Adds the provided rikishi to the database
func (s *RikishiService) AddRikishi(rikishi []models.Rikishi) error {
	s.log.Infof("adding %d rikishi", len(rikishi))
	var err error
	var ranks map[string]models.Rank
	if ranks, err = s.rankService.GetRanksMappedByName(); err != nil {
		s.log.Errorf("error while retriving ranks: %s", err)
		return err
	}
	count := 0
	for index, curRikishi := range rikishi {
		if rank, ok := ranks[curRikishi.Rank]; !ok {
			s.log.Infof("unknown rikshi rank %s for riksihi #%d", curRikishi.Rank, index)
		} else {
			s.log.Infof("binding rikishi '%s' to rank '%s'", curRikishi.Name, rank.Name)
			s.db.Model(&rank).Association("Rikishis").Append(&curRikishi)
			count++
		}
	}

	for _, rank := range ranks {
		s.db.Save(&rank)
	}

	s.log.Infof("successfully added %d rikishi", count)
	return nil
}
