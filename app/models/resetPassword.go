package models

import (
	"time"
)

type ResetPassword struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ResetToken string `gorm:"size:255;not null;uniqueIndex" json:"reset_token"`
	IPAddress  string `gorm:"size:255;not null" json:"ip_address" validate:"required"`
	Device     string `gorm:"size:255;not null" json:"device" validate:"required"`
	EmailSent  int    `gorm:"size:1;not null" json:"email_sent"`
	UUID       string `gorm:"size:36;not null" json:"resetpw_uuid"`
	User       User   `gorm:"foreignKey:UUID;references:UUID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
