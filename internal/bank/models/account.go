package models

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Account struct {
	ID            uuid.UUID      `json:"id" gorm:"primary_key"`
	AccountNumber ulid.ULID      `json:"account_number" gorm:"not null"`
	UserId        uuid.UUID      `json:"user_id" gorm:"not null"`
	Name          string         `json:"name" gorm:"not null"`
	Balance       int64          `json:"balance" gorm:"not null"`
	CurrencyId    int            `json:"currency_id" gorm:"not null"`
	CreatedAt     *time.Time     `gorm:":autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:":autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"default:null"`
}

func (a *Account) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New()
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	a.AccountNumber = ulid.MustNew(ulid.Timestamp(time.Now()), ulid.Monotonic(entropy, 0))
	return nil
}
