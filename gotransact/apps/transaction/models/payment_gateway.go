package models

import (
	base "gotransact/apps/base/models"

	"gorm.io/gorm"
)

type PaymentGateway struct {
	gorm.Model
	base.Base
	Slug               string             `gorm:"type:varchar(100);not null;unique"`
	Label              string             `gorm:"type:varchar(100);not null"`
	TransactionRequest TransactionRequest `gorm:"foreignKey:PaymentGatewayMethodID"`
}
