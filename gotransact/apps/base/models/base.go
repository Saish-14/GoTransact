package models

import (
	"github.com/google/uuid"
)

type Base struct {
	InternalID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Is_Active  bool      `gorm:"default:true"`
}
