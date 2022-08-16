package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	ID       uint   `json:"id"`
	Content  string `gorm:"type:text" json:"content"`
	AuthorID uint
	Author   User `gorm:"foreignKey:AuthorID"`
}
