package services

import (
	"reflect"

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
	matchService    MatchService

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
	matchService MatchService,
) LoadService {
	return LoadService{
		log:             log,
		db:              db,
		data:            data,
		rankService:     rankService,
		rikishiService:  rikishiService,
		categoryService: categoryService,
		matchService:    matchService,
	}
}

// Load - Loads data from the data file into the given table
func (s LoadService) Load(model string, clear bool) error {
	s.log.Infof("Loading '%s'", model)
	switch model {
	case "category":
		slice := s.parseModel(&models.Category{})
		return s.categoryService.SetCategories(*(slice.(*[]models.Category)))
	case "rikishi":
		if clear {
			s.clearModel(&models.Rikishi{})
		}
		slice := s.parseModel(&models.Rikishi{})
		return s.rikishiService.AddRikishi(*(slice.(*[]models.Rikishi)))
	case "match":
		if clear {
			s.clearModel(&models.Match{})
			slice := s.parseModel(&models.Match{})
			return s.matchService.AddMatches(*(slice.(*[]models.Match)))
		}
	}
	var obj interface{}
	var err error
	if obj, err = models.GetModelFromID(model); err != nil {
		s.log.Errorf("error while loading model '%s': %s", model, err)
	}
	return s.loadModel(obj, clear)
}

func (s LoadService) clearModel(model interface{}) {
	modelType := reflect.ValueOf(model).Elem().Type()

	s.log.Infof("clearing '%s' from the database", modelType.Name())
	s.db.Unscoped().Where("1 = 1").Delete(model)
}

func (s LoadService) parseModel(model interface{}) interface{} {
	modelType := reflect.ValueOf(model).Elem().Type()
	s.log.Infof("parsing model %s", modelType)
	modelSlice := reflect.New(reflect.SliceOf(modelType))
	s.data.UnmarshalKey(modelType.Name(), modelSlice.Interface())

	for index := 0; index < modelSlice.Elem().Len(); index++ {
		s.log.Infof("[%d]=%v", index, modelSlice.Elem().Index(index))
	}
	return modelSlice.Interface()
}

func (s LoadService) uploadData(model interface{}, data interface{}) error {
	modelType := reflect.ValueOf(model).Elem().Type()
	if err := s.db.Create(data).Error; err != nil {
		s.log.Errorf("error loading %s into the database: %s", modelType.Name(), err)
		return err
	}

	s.log.Infof("loaded %s", modelType.Name())
	return nil
}

func (s LoadService) loadModel(model interface{}, clear bool) error {
	modelType := reflect.ValueOf(model).Elem().Type()
	s.log.Infof("loading '%s' from config", modelType.Name())

	if clear {
		s.log.Infof("clearing '%s' from the database", modelType.Name())
		s.db.Unscoped().Where("1 = 1").Delete(model)
	}

	modelSlice := reflect.New(reflect.SliceOf(modelType))
	s.data.UnmarshalKey(modelType.Name(), modelSlice.Interface())

	for index := 0; index < modelSlice.Elem().Len(); index++ {
		s.log.Infof("[%d]=%v", index, modelSlice.Elem().Index(index))
	}

	if err := s.db.Create(modelSlice.Interface()).Error; err != nil {
		s.log.Errorf("error loading %s into the database: %s", modelType.Name(), err)
		return err
	}

	s.log.Infof("loaded %s", modelType.Name())
	return nil
}
