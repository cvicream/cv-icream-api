package model

import (
	"time"

	"gorm.io/gorm"
)

type Survey struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;constraint:OnDelete:CASCADE" json:"userId"`
	Survey      string         `gorm:"size:255;not null" json:"survey"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
