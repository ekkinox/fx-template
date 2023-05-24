package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
}
