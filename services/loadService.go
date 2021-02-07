package services

import (
	"fmt"

	"github.com/Jeff-All/Dohyo/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// LoadService - Service used to load data from data files into the database
type LoadService struct {
	log *logrus.Logger
	db  *gorm.DB

	rankService    RankService
	rikishiService RikishiService

	data *viper.Viper
}

// NewLoadService - Returns a new LoadService
func NewLoadService(
	log *logrus.Logger,
	db *gorm.DB,
	data *viper.Viper,
	rankService RankService,
	rikishiService RikishiService,
) LoadService {
	return LoadService{
		log:            log,
		db:             db,
		data:           data,
		rankService:    rankService,
		rikishiService: rikishiService,
	}
}

// Load - Loads data from the data file into the given table
func (s LoadService) Load(model string) error {
	s.log.Infof("Loading '%s'", model)
	switch model {
	case "rank":
		return s.LoadRanks()
	case "rikishi":
		return s.LoadRikishi()
	default:
		s.log.Errorf("invalid model '%s'", model)
		return fmt.Errorf("invalid model '%s'", model)
	}
}

// LoadRanks - Fills the Ranks table with data from the data file
func (s LoadService) LoadRanks() error {
	s.log.Info("loading ranks from config")
	var ranks []models.Rank
	s.data.UnmarshalKey("ranks", &ranks)

	s.log.Infof("loading %d rank entries", len(ranks))

	fmt.Println(s.data.GetStringMap("ranks"))

	s.db.Unscoped().Exec("DELETE FROM ranks")

	s.db.Create(&ranks)
	s.log.Info("loaded ranks")
	return nil
}

// LoadRikishi - Fills the Ranks table with data from the data file
func (s LoadService) LoadRikishi() error {
	s.log.Info("loading rikishi from config")
	var rikishi []models.Rikishi
	s.data.UnmarshalKey("rikishi", &rikishi)

	s.rikishiService.AddRikishi(rikishi)

	s.log.Info("loaded rikishi")
	return nil
}
