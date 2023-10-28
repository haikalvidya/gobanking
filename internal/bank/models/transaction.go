package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	WITHDRAWAL = "withdrawal"
	DEPOSIT    = "deposit"
	TRANSFER   = "transfer"
)

type Transaction struct {
	ID              uuid.UUID      `json:"id" gorm:"primary_key"`
	SourceAccountId uuid.UUID      `json:"source_account_id" gorm:"not null"`
	DestAccountId   uuid.UUID      `json:"dest_account_id"`
	Amount          int64          `json:"amount" gorm:"not null"`
	Type            string         `json:"type" gorm:"not null"`
	CreatedAt       *time.Time     `gorm:":autoCreateTime" json:"created_at"`
	DeletedAt       gorm.DeletedAt `gorm:"default:null"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}
