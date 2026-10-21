package models

import "gorm.io/gorm"

type ArticleLike struct {
	gorm.Model
	UserID    uint `gorm:"uniqueIndex:idx_user_article"`
	ArticleID uint `gorm:"uniqueIndex:idx_user_article"`
}
