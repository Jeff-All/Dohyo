package services

import (
	"errors"
	"strconv"

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

// GetRikishiByCategory - Returns a map of all rikishi indexed by their category
func (s *CategoryService) GetRikishiByCategory() (map[string][]string, error) {
	s.log.Infof("retrieving rikishi by their categories")

	var categories []models.Category
	var err error
	if categories, err = s.GetAllCategoriesWithRikishis(); err != nil {
		s.log.Errorf("eror while retrieving categories: %s", err)
		return nil, err
	}

	rikishiMap := make(map[string][]string)
	for _, category := range categories {
		array := make([]string, len(category.Rikishis))
		for index, rikishi := range category.Rikishis {
			array[index] = strconv.FormatUint(uint64(rikishi.ID), 10)
		}
		rikishiMap[category.Name] = array
	}

	return rikishiMap, nil
}

// Count - Returns the number of categories currently configured in the DB
func (s *CategoryService) Count() (int, error) {
	var count int64 = 0
	err := s.db.Model(&models.Category{}).Count(&count).Error
	if err != nil {
		s.log.Errorf("error pulling count from categories: %s", err)
		return 0, err
	}
	return int(count), nil
}

// GetCategoryCountOfRikishis - Gets the number of distinct categories of the rikishis
func (s *CategoryService) GetCategoryCountOfRikishis(rikishis models.Rikishis) (int, error) {
	rikishiIDs := rikishis.GetIDs()
	var count int64 = 0
	var err error
	if err = s.db.Table("rikishis").Select("count(distinct(categories.id))").Joins("JOIN categories ON categories.id = rikishis.category_id").Where("rikishis.id IN ?", rikishiIDs).Count(&count).Error; err != nil {
		s.log.Errorf("error pulling category count: %s", err)
		return 0, err
	}
	return int(count), nil
}

type rikishiCategory struct {
	models.Rikishi
	Category string
}

// GetRikishisByCategoryByID - Returns the provided rikishis mapped to by their categories
func (s *CategoryService) GetRikishisByCategoryByID(rikishis models.Rikishis) (map[string]models.Rikishi, error) {
	rikishiIDs := rikishis.GetIDs()
	rikishiCategories := make([]rikishiCategory, 0, len(rikishis))
	if err := s.db.Table("rikishis").Select("rikishis.*, categories.name AS category").Joins("JOIN categories ON categories.id = rikishis.category_id").Where("rikishis.id IN ?", rikishiIDs).Find(&rikishiCategories).Error; err != nil {
		s.log.Errorf("error getting rikishi categories: %s", err)
		return nil, err
	}
	rikishiMap := make(map[string]models.Rikishi)
	for _, rikishiCategory := range rikishiCategories {
		rikishiMap[rikishiCategory.Category] = rikishiCategory.Rikishi
	}
	return rikishiMap, nil
}
