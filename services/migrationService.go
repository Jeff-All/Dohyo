package services

import (
	"reflect"

	"github.com/Jeff-All/Dohyo/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// MigrationService - Service for performing migrations
type MigrationService struct {
	log *logrus.Logger
	db  *gorm.DB
}

// NewMigrationService - Instantiates a new MigrationService with the given parameters
func NewMigrationService(log *logrus.Logger, db *gorm.DB) MigrationService {
	return MigrationService{
		log: log,
		db:  db,
	}
}

// MigrateModels - Migrates all the models supplied
func (s *MigrationService) MigrateModels(ids ...string) error {
	s.log.Info("migrating models", ids)
	objects := make([]interface{}, len(ids))
	for index, id := range ids {
		var model interface{}
		var err error
		if model, err = models.GetModelFromID(id); err != nil {
			return err
		}
		reflect.ValueOf(model)
		objects[index] = model
		_ = append(objects, model)
	}
	s.db.AutoMigrate(objects...)
	return nil
}
