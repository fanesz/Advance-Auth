package config

import (
	"gorm.io/gorm"
	"advanceauth/backend/app/models"
)

type MigratableModel interface {
	Migrate(db *gorm.DB) error
}

var listModels = []interface{}{
	&models.User{},
	&models.LoggedUser{},
	&models.ResetPassword{},
}

func Migrate(db *gorm.DB) error {
	for _, model := range listModels {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}
	return nil
}
