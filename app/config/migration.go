package config

import (
	"advanceauth/backend/app/models"
	"gorm.io/gorm"
)

type MigratableModel interface {
	Migrate(db *gorm.DB) error
}

var listModels = []interface{}{
	&models.User{},
	&models.LoggedUser{},
	&models.ResetPassword{},
	&models.VerifyUser{},
}

func Migrate(db *gorm.DB) error {
	for _, model := range listModels {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}
	return nil
}
