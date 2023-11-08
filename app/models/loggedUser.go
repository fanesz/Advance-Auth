package models

import (
	"gorm.io/gorm"
	"time"
)

type LoggedUser struct {
	ID         uint   `gorm:"primary_key;auto_increment" json:"id"`
	LoginToken string `gorm:"size:255;not null" json:"login_token"`
	IPAddress  string `gorm:"size:255;not null" json:"ip_address" validate:"required"`
	Device     string `gorm:"size:255;not null" json:"device" validate:"required"`
	UUID       string `gorm:"size:36;not null" json:"logged_uuid"`
	User       User   `gorm:"foreignKey:UUID;references:UUID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func MigrateLoggedUser(db *gorm.DB) error {
	if err := db.AutoMigrate(&LoggedUser{}); err != nil {
		return err
	}
	return nil
}
