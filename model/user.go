package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uint           `json:"id"`
	Role        string         `gorm:"not null" json:"type"`
	Cpf         string         `gorm:"not null,unique" json:"cpf"`
	Name        string         `gorm:"not null,type:text" json:"name"`
	Password    []byte         `gorm:"not null" json:"password"`
	Email       string         `gorm:"not null,unique" json:"email"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index,->" json:"-"`
}

type UserToken struct {
	Cpf         string `json:"cpf"`
	Role        string `json:"role"`
	TokenString string `json:"token"`
}
