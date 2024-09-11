package model

import (
	"time"

	"gorm.io/gorm"
)

type CV struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null" json:"userId"`
	Title       *string        `gorm:"size:100" json:"title"`
	Description *string        `gorm:"size:255" json:"description"`
	Content     *string        `gorm:"type:text" json:"content"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	User        User
}
