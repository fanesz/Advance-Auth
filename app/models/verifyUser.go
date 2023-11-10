package models

import (
	"time"
)

type VerifyUser struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	VerifyToken string `gorm:"size:255;not null;uniqueIndex" json:"verify_token"`
	EmailSent   int    `gorm:"size:1;not null" json:"email_sent"`
	UUID        string `gorm:"size:36;not null" json:"verifying_uuid"`
	User        User   `gorm:"foreignKey:UUID;references:UUID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
