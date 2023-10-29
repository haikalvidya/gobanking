package models

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Account struct {
	// account number is ulid but the type is uuid.UUID because gorm doesn't support ulid
	AccountNumber string         `json:"account_number" gorm:"type:varchar(36);not null;primary_key;"`
	UserId        uuid.UUID      `json:"user_id" gorm:"not null"`
	Name          string         `json:"name" gorm:"not null"`
	Balance       int64          `json:"balance" gorm:"not null"`
	CurrencyId    int            `json:"currency_id" gorm:"not null"`
	CreatedAt     *time.Time     `gorm:":autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:":autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"default:null"`
}

func (a *Account) BeforeCreate(tx *gorm.DB) error {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	a.AccountNumber = ulid.MustNew(ulid.Timestamp(time.Now()), ulid.Monotonic(entropy, 0)).String()
	return nil
}
