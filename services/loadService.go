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

	rankService     RankService
	rikishiService  RikishiService
	categoryService CategoryService

	data *viper.Viper
}

// NewLoadService - Returns a new LoadService
func NewLoadService(
	log *logrus.Logger,
	db *gorm.DB,
	data *viper.Viper,
	rankService RankService,
	rikishiService RikishiService,
	categoryService CategoryService,
) LoadService {
	return LoadService{
		log:             log,
		db:              db,
		data:            data,
		rankService:     rankService,
		rikishiService:  rikishiService,
		categoryService: categoryService,
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
	case "category":
		return s.LoadCategories()
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

	if err := s.rikishiService.AddRikishi(rikishi); err != nil {
		s.log.Errorf("error while adding rikishi: %s", err)
		return err
	}

	s.log.Info("loaded rikishi")
	return nil
}

// LoadCategories - Fills the Categories tables with data from the data file
func (s LoadService) LoadCategories() error {
	s.log.Info("loading categories from config")

	var categories []models.Category
	s.data.UnmarshalKey("categories", &categories)

	if err := s.categoryService.SetCategories(categories); err != nil {
		s.log.Errorf("error while adding categories: %s", err)
		return err
	}

	s.log.Info("loaded categories")
	return nil
}
