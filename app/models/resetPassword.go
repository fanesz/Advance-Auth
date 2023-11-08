package models

import (
	"gorm.io/gorm"
	"time"
)

type ResetPassword struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ResetToken string `gorm:"size:255;not null" json:"reset_token"`
	IPAddress  string `gorm:"size:255;not null" json:"ip_address" validate:"required"`
	Device     string `gorm:"size:255;not null" json:"device" validate:"required"`
	EmailSent  int    `gorm:"size:1;not null" json:"email_sent"`
	UUID       string `gorm:"size:36;not null" json:"logged_uuid"`
	User       User   `gorm:"foreignKey:UUID;references:UUID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func MigrateResetPassword(db *gorm.DB) error {
	if err := db.AutoMigrate(&ResetPassword{}); err != nil {
		return err
	}
	return nil
}
