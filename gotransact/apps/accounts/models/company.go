package models

import (
	base "gotransact/apps/base/models"

	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	base.Base
	Name   string `gorm:"type:varchar(100);not null"`
	UserID uint   `gorm:"not null"`
}
