package dto

import (
	"time"
)

type CreateAccountRequest struct {
	Name         string `json:"name" validate:"required"`
	CurrencyCode string `json:"currency_code" validate:"required"`
}

type UpdateAccountRequest struct {
	Name string `json:"name"`
}

type GetAccountByExternalResponse struct {
	AccountNumber string     `json:"account_number" gorm:"type:varchar(36);not null;primary_key;"`
	Name          string     `json:"name" gorm:"not null"`
	CurrencyId    int        `json:"currency_id" gorm:"not null"`
	CreatedAt     *time.Time `gorm:":autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time `gorm:":autoUpdateTime" json:"updated_at"`
}
