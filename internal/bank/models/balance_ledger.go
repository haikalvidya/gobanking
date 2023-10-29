package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BalanceLedger struct {
	ID            uuid.UUID      `json:"id" gorm:"primary_key"`
	AccountNumber string         `json:"account_number" gorm:"type:varchar(36);not null;"`
	Balance       int64          `json:"balance" gorm:"not null"`
	TransactionId uuid.UUID      `json:"transaction_id" gorm:"not null"`
	CreatedAt     *time.Time     `gorm:":autoCreateTime" json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"default:null"`
}

func (b *BalanceLedger) BeforeCreate(tx *gorm.DB) error {
	b.ID = uuid.New()
	return nil
}
