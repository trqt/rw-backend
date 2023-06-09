package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content  string `gorm:"type:text" json:"content"`
	Rating   uint   `json:"rating"`
	AuthorID uint   `json:"author_id"`
	WorkerID uint   `json:"worker_id"`
	GigID    uint   `json:"gig_id"`
}
