package models

import (
	base "gotransact/apps/base/models"

	"gorm.io/gorm"
)


type TransactionHistory struct {
	gorm.Model
	base.Base
	TransactionID uint              `json:"transactionid" gorm:"not null"`
	Status        TransactionStatus `json:"status" gorm:"type:varchar(20);not null"`
	Description   string            `json:"description" gorm:"size:255"`
	Amount        float64           `json:"amount" gorm:"type:float"`
}
