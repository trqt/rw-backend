package model

import "gorm.io/gorm"

type Gig struct {
	gorm.Model
	ID        uint   `json:"ID"`
	Completed bool   `gorm:"not null" json:"completed"`
	Approved  bool   `gorm:"not null" json:"approved"`
	Desc      string `gorm:"not null" json:"description"`

	WorkerID uint `json:"worker_id"`
	//Worker   User `gorm:"foreignKey:WorkerID"`

	HirerID uint `json:"hirer_id"`
	//Hirer   User `gorm:"foreignKey:HirerID"`
}
