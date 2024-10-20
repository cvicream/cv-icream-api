package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	FirstName  *string        `gorm:"size:50" json:"firstName"`
	LastName   *string        `gorm:"size:50" json:"lastName"`
	Email      string         `gorm:"size:255;not null" json:"email"`
	Password   *string        `gorm:"size:255" json:"-"`
	Avatar     string         `gorm:"type:text" json:"avatar"`
	Provider   string         `gorm:"size:50;not null" json:"provider"`
	ProviderID *string        `gorm:"size:50" json:"providerId"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
