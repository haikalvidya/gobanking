package repository

import (
	"context"
	"gobanking/internal/bank/models"

	"gorm.io/gorm"
)

type AccountRepository interface {
	// Create creates a new account
	CreateTX(ctx context.Context, tx *gorm.DB, account *models.Account) error
	// GetByAccountNumber returns an account by account number
	GetByAccountNumber(ctx context.Context, accountNumber string) (*models.Account, error)
	// GetByUserID returns all account by user id
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.Account, error)
	// GetByAccountNumberAndUserID returns an account by account number and user id
	GetByAccountNumberAndUserID(ctx context.Context, accountNumber string, userID string) (*models.Account, error)
	// Update updates an account
	UpdateTX(ctx context.Context, tx *gorm.DB, account *models.Account) error
	// Delete deletes an account
	Delete(ctx context.Context, accountNumber string) error
}

type accountRepository repository

func (r *accountRepository) CreateTX(ctx context.Context, tx *gorm.DB, account *models.Account) error {
	return tx.WithContext(ctx).Create(account).Error
}

func (r *accountRepository) GetByAccountNumber(ctx context.Context, accountNumber string) (*models.Account, error) {
	account := &models.Account{}
	err := r.DB.WithContext(ctx).Where("account_number = ?", accountNumber).First(account).Error
	return account, err
}

func (r *accountRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.Account, error) {
	var accounts []*models.Account
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) GetByAccountNumberAndUserID(ctx context.Context, accountNumber string, userID string) (*models.Account, error) {
	account := &models.Account{}
	err := r.DB.WithContext(ctx).Where("account_number = ? AND user_id = ?", accountNumber, userID).First(account).Error
	return account, err
}

func (r *accountRepository) UpdateTX(ctx context.Context, tx *gorm.DB, account *models.Account) error {
	return tx.WithContext(ctx).Save(account).Error
}

func (r *accountRepository) Delete(ctx context.Context, accountNumber string) error {
	return r.DB.WithContext(ctx).Where("account_number = ?", accountNumber).Delete(&models.Account{}).Error
}
