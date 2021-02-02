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

	data *viper.Viper
}

// NewLoadService - Returns a new LoadService
func NewLoadService(
	log *logrus.Logger,
	db *gorm.DB,
	data *viper.Viper,
) LoadService {
	return LoadService{
		log:  log,
		db:   db,
		data: data,
	}
}

// Load - Loads data from the data file into the given table
func (s LoadService) Load(model string) {
	s.log.Infof("Loading '%s'", model)
	switch model {
	case "rank":
		s.LoadRanks()
	case "rikishi":
		s.LoadRikishi()
	}
}

// LoadRanks - Fills the Ranks table with data from the data file
func (s LoadService) LoadRanks() {
	s.log.Info("populating ranks from config")
	var ranks []models.Rank
	s.data.UnmarshalKey("ranks", &ranks)

	fmt.Println(s.data.GetStringMap("ranks"))

	s.db.Unscoped().Exec("DELETE FROM ranks")

	s.db.Create(&ranks)
	s.log.Info("populated ranks")
}

// LoadRikishi - Fills the Ranks table with data from the data file
func (s LoadService) LoadRikishi() {
	s.log.Info("populating rikishi from config")
	var rikishi []models.Rikishi
	s.data.UnmarshalKey("rikishi", &rikishi)

	s.db.Create(&rikishi)
	s.log.Info("populated rikishi")
}
