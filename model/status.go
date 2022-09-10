package model

import "gorm.io/gorm"

type WorkStatus struct {
	gorm.Model
	ID     uint   `json:"id"`
	Status string `gorm:"type:text" json:"status"`

	WorkerID uint `json:"worker_id"`
	Worker   User `gorm:"foreignKey:WorkerID"`

	HirerID uint `json:"hirer_id"`
	Hirer   User `gorm:"foreignKey:HirerID"`
}
