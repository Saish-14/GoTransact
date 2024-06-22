package models

import (
	base "gotransact/apps/base/models"

	"gorm.io/gorm"
)

type TransactionStatus string

const (
	StatusPending    TransactionStatus = "pending"
	StatusProcessing TransactionStatus = "processing"
	StatusSuccess    TransactionStatus = "success"
	StatusFailed     TransactionStatus = "failed"
)

type TransactionRequest struct {
	gorm.Model
	base.Base
	UserID uint `gorm:""`
	Status             TransactionStatus  `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	PaymentGatewayMethodID uint   `gorm:"not null"`
	Description        string             `gorm:"type:text"`
	Amount             float64             `gorm:"type:string;not null"`
	TransactionHistory TransactionHistory `gorm:"foreignKey:TransactionID"`
}

