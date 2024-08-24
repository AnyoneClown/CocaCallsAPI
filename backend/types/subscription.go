package types

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	Plan      string
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

