package model

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID          uint           `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt   time.Time      `json:"created_at" form:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" form:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
	Title       string         `json:"title" form:"title"`
	Description string         `json:"description" form:"description"`
	Likes       int            `json:"likes" form:"likes"`
}
