package services

import (
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
