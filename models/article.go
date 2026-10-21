package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title     string `gorm:"size:200;index" json:"title"`
	Content   string `gorm:"type:text" json:"content"`
	UserID    uint   `json:"user_id"`
	Author    User   `gorm:"foreignKey:UserID" json:"author"`
	LikeCount int    `json:"like_count"`
}
