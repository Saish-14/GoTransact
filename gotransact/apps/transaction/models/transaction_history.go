package models

import (
	base "gotransact/apps/base/models"

	"gorm.io/gorm"
)

type TransactionHistory struct {
	gorm.Model
	base.Base
	TransactionID uint `gorm:"not null"`
	Status      string `gorm:"type:varchar(20);not null"`
	Description string `gorm:"type:text"`
	Amount      string `gorm:"type:string;not null"`
}

// // Enum for status
// const (
// 	StatusPending    = "pending"
// 	StatusProcessing = "processing"
// 	StatusSuccess    = "success"
// 	StatusFailed     = "failed"
// )
