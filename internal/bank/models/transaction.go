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
	ID                  uuid.UUID      `json:"id" gorm:"primary_key"`
	SourceAccountNumber string         `json:"source_account_number" gorm:"type:varchar(36)"`
	DestAccountNumber   string         `json:"dest_account_number" gorm:"type:varchar(36)"`
	Amount              int64          `json:"amount" gorm:"not null"`
	Type                string         `json:"type" gorm:"not null"`
	CreatedAt           *time.Time     `gorm:":autoCreateTime" json:"created_at"`
	DeletedAt           gorm.DeletedAt `gorm:"default:null"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}
