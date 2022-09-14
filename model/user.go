package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uint           `json:"ID"`
	Role        string         `gorm:"not null" json:"role"`
	Cpf         string         `gorm:"not null,unique" json:"cpf,omitempty"`
	Name        string         `gorm:"not null,type:text" json:"name"`
	Password    string         `gorm:"not null" json:"password,omitempty"`
	Email       string         `gorm:"not null,unique" json:"email,omitempty"`
	Description string         `gorm:"type:text" json:"description"`
	Tags        string         `gorm:"not null,type:text" json:"tags"`
	Photo       string         `gorm:"not null" json:"photo_url"`
	Phone       string         `gorm:"not null" json:"phone_number"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index,->" json:"-"`
}
