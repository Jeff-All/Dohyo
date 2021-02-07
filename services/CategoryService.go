package services

import (
	"errors"

	"github.com/Jeff-All/Dohyo/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CategoryService - Service for handling categories
type CategoryService struct {
	db  *gorm.DB
	log *logrus.Logger

	rikishiService RikishiService
}

// NewCategoryService - Instantiates a new CategoryService
func NewCategoryService(
	log *logrus.Logger,
	db *gorm.DB,
	rikishiService RikishiService,
) CategoryService {
	return CategoryService{
		db:             db,
		log:            log,
		rikishiService: rikishiService,
	}
}

// SetCategories - Overwrites categories with the provided categories
func (s *CategoryService) SetCategories(categories []models.Category) error {
	var rikishis map[string]models.Rikishi
	var err error
	if rikishis, err = s.rikishiService.GetRikishiMappedByName(); err != nil {
		s.log.Errorf("error while pulling rikishi entries: %s", err)
	}

	s.db.Unscoped().Exec("DELETE FROM categories")

	s.db.Create(categories)

	for _, category := range categories {
		for _, rikishiName := range category.RikishiNames {
			if rikishi, ok := rikishis[rikishiName]; ok {
				s.log.Infof("binding rikishi '%s' to category '%s'", rikishi.Name, category.Name)
				s.db.Model(&category).Association("Rikishis").Append(&rikishi)
			} else {
				s.log.Warnf("unable to locate rikishi '%s' for category '%s'", rikishiName, category.Name)
			}
		}
	}

	for category := range categories {
		s.db.Save(category)
	}

	return nil
}

// GetAllCategories - Returns all the categories in the Database
func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	s.log.Info("retrieving all categories")
	categories := []models.Category{}
	var result *gorm.DB
	if result = s.db.Find(&categories); result.Error == nil {
		s.log.Infof("successfully pulled all categories")
		return categories, nil
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		s.log.Infof("unable to find any category entries")
		return categories, nil
	}
	s.log.Errorf("error while retrieving category entries: %s", result.Error)
	return categories, result.Error
}

// GetAllCategoriesWithRikishis - Returns all the categories in the DB with their
// associated rikishi
func (s *CategoryService) GetAllCategoriesWithRikishis() ([]models.Category, error) {
	s.log.Info("retrieving all categories")
	categories := []models.Category{}
	var result *gorm.DB
	if result = s.db.Preload("Rikishis").Find(&categories); result.Error == nil {
		s.log.Infof("successfully pulled all categories")
		return categories, nil
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		s.log.Infof("unable to find any category entries")
		return categories, nil
	}
	s.log.Errorf("error while retrieving category entries: %s", result.Error)
	return categories, result.Error
}

// GetRikishiByCategory - Returns a map of rikishi indexed b their category
func (s *CategoryService) GetRikishiByCategory() (map[string][]models.Rikishi, error) {
	s.log.Infof("retrieving rikishi by their categories")

	var categories []models.Category
	var err error
	if categories, err = s.GetAllCategoriesWithRikishis(); err != nil {
		s.log.Errorf("eror while retrieving categories: %s", err)
		return nil, err
	}

	rikishiMap := make(map[string][]models.Rikishi)
	for _, category := range categories {
		array := make([]models.Rikishi, len(category.Rikishis))
		for index, rikishi := range category.Rikishis {
			array[index] = rikishi
		}
		rikishiMap[category.Name] = array
	}

	return rikishiMap, nil
}
