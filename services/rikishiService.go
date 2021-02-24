package services

import (
	"fmt"

	"github.com/Jeff-All/Dohyo/models"
	"github.com/Jeff-All/Dohyo/responses"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RikishiService - Provides functions for accessing the Rikishi data
type RikishiService struct {
	log               *logrus.Logger
	db                *gorm.DB
	rankService       RankService
	TournamentService TournamentService
}

// NewRikishiService - Instantiates a new RikishiService
func NewRikishiService(
	log *logrus.Logger,
	db *gorm.DB,
	rankService RankService,
	tournamentService TournamentService,
) RikishiService {
	return RikishiService{
		log:               log,
		db:                db,
		rankService:       rankService,
		TournamentService: tournamentService,
	}
}

// GetAllCurrentCompleteRikishi - returns all the rikishi with their current tournament's matches
func (s RikishiService) GetAllCurrentCompleteRikishi() ([]responses.Rikishi, error) {
	rikishis := []responses.Rikishi{}
	if err := s.db.Raw("SELECT id, name, avatar, rank FROM rikishis_complete ORDER BY id").Scan(&rikishis).Error; err != nil {
		s.log.Errorf("error while pulling all complete rikishi: %s", err)
		return nil, err
	}
	var tournament *models.Tournament
	var err error
	if tournament, err = s.TournamentService.GetCurrentTournament(); err != nil {
		s.log.Errorf("error getting current tournament: %s", err)
	} else if tournament.ID != 0 {
		s.log.Infof("scanning current matches into rikishis")
		for i := 0; i < len(rikishis); i++ {
			rikishi := &rikishis[i]
			if err = s.PopulateTournamentMatchesForRikishi(rikishi, *tournament); err != nil {
				s.log.Errorf("error populating tournament matches: %s", err)
				return nil, err
			}
			if err = s.PopulateTournamentResultsForRikishi(rikishi, 5); err != nil {
				s.log.Errorf("error populating tournament results: %s", err)
				return nil, err
			}
		}
	} else {
		s.log.Infof("there is no current tournament set")
	}

	return rikishis, nil
}

// PopulateTournamentMatchesForRikishi - Fills matches for the given rikishi
func (s *RikishiService) PopulateTournamentMatchesForRikishi(rikishi *responses.Rikishi, tournament models.Tournament) error {
	s.log.Infof("scanning matches into '%d'", rikishi.ID)
	matches := []responses.Match{}
	if err := s.db.Raw("SELECT day, opponent, concluded, won FROM rikishi_matches WHERE tournament_id = ? AND rikishi_id = ? ORDER BY day", tournament.ID, rikishi.ID).Scan(&matches).Error; err != nil {
		s.log.Errorf("error loading matches for rikishi '%d': %s", rikishi.ID, err)
		return err
	}
	rikishi.Matches = make(map[uint]responses.Match, len(matches))
	for _, match := range matches {
		rikishi.Matches[match.Day] = match
		if match.Concluded {
			if match.Won {
				rikishi.Wins++
			} else {
				rikishi.Losses++
			}
		}
	}
	return nil
}

// PopulateTournamentResultsForRikishi - Fills matches for the given rikishi
func (s *RikishiService) PopulateTournamentResultsForRikishi(
	rikishi *responses.Rikishi,
	count int,
) error {
	s.log.Infof("scanning results into '%d'", rikishi.ID)
	results := []responses.Result{}
	if err := s.db.Raw("SELECT tournament, wins, losses FROM tournament_results WHERE rikishi_id = ? ORDER BY tournament_id ASC LIMIT ?", rikishi.ID, count).Scan(&results).Error; err != nil {
		s.log.Errorf("error loading results for rikishi '%d': %s", rikishi.ID, err)
		return err
	}
	rikishi.Results = results
	return nil
}

// GetAllRikishi - Returns all rikishi in the database
func (s *RikishiService) GetAllRikishi() ([]models.Rikishi, error) {
	ranks := make([]models.Rank, 0)
	if err := s.db.Preload("Rikishis").Find(&ranks).Error; err != nil {
		s.log.Errorf("error while pulling ranks and rikishis: %s", err)
		return nil, err
	}

	rikishis := []models.Rikishi{}
	for _, rank := range ranks {
		for _, rikishi := range rank.Rikishis {
			rikishi.Rank = rank.Name
			if rikishi.SubRank > 0 {
				rikishi.Rank += fmt.Sprintf(" %d", rikishi.SubRank)
			}
			rikishis = append(rikishis, rikishi)
		}
	}
	s.log.Infof("successfully pulled all rikishi")
	return rikishis, nil
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
