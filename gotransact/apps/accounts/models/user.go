package models


import (
	base "gotransact/apps/base/models"
	transaction "gotransact/apps/transaction/models"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	base.Base
	FirstName          string                         `gorm:"type:varchar(100);not null"`
	LastName           string                         `gorm:"type:varchar(100);not null"`
	Email              string                         `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password           string                         `gorm:"type:varchar(255)"`
	Company            Company                        `gorm:"foreignKey:UserID"`
	TransactionRequest transaction.TransactionRequest `gorm:"foreignKey:UserID"`
}
