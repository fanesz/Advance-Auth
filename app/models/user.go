package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID     string `gorm:"size:36;not null;uniqueIndex;primary_key" json:"uuid"`
	Username string `gorm:"size:255;not null;uniqueIndex" json:"username"`
	Email    string `gorm:"size:255;not null;uniqueIndex" json:"email" validate:"required,email"`
	Password string `gorm:"size:255;not null" json:"password" validate:"required"`
}

func MigrateUsers(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}
