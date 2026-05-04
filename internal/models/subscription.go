package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ServiceName string         `gorm:"not null" json:"service_name"`
	Price       int            `gorm:"not null" json:"price"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	StartDate   time.Time      `gorm:"not null" json:"start_date"`
	EndDate     *time.Time     `json:"end_date,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}
