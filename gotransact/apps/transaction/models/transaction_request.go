package models

import (
	base "gotransact/apps/base/models"

	"gorm.io/gorm"
)

type TransactionRequest struct {
	gorm.Model
	base.Base
	UserID uint `gorm:""`
	Status                 string `gorm:"type:varchar(20);not null"`
	PaymentGatewayMethodID uint   `gorm:"not null"`
	Description        string             `gorm:"type:text"`
	Amount             string             `gorm:"type:string;not null"`
	TransactionHistory TransactionHistory `gorm:"foreignKey:TransactionID"`
}

// Enum for status
const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusSuccess    = "success"
	StatusFailed     = "failed"
)
