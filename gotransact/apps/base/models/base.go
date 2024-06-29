package models

import (
	"github.com/google/uuid"
)

type Base struct {
	InternalID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Is_Active  bool      `gorm:"default:true"`
}
