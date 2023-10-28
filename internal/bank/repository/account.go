package repository

import (
	"context"
	"gobanking/internal/bank/models"
)

type AccountRepository interface {
	// Create creates a new account
	Create(ctx context.Context, account *models.Account) error
	// GetByAccountNumber returns an account by account number
	GetByAccountNumber(ctx context.Context, accountNumber string) (*models.Account, error)
	// GetByUserID returns all account by user id
	GetByUserID(ctx context.Context, userID string) ([]*models.Account, error)
	// GetByID returns an account by id
	GetByID(ctx context.Context, id string) (*models.Account, error)
	// Update updates an account
	Update(ctx context.Context, account *models.Account) error
	// Delete deletes an account
	Delete(ctx context.Context, id string) error
}

type accountRepository repository

func (r *accountRepository) Create(ctx context.Context, account *models.Account) error {
	return r.DB.WithContext(ctx).Create(account).Error
}

func (r *accountRepository) GetByAccountNumber(ctx context.Context, accountNumber string) (*models.Account, error) {
	account := &models.Account{}
	err := r.DB.WithContext(ctx).Where("account_number = ?", accountNumber).First(account).Error
	return account, err
}

func (r *accountRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Account, error) {
	var accounts []*models.Account
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) GetByID(ctx context.Context, id string) (*models.Account, error) {
	account := &models.Account{}
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(account).Error
	return account, err
}

func (r *accountRepository) Update(ctx context.Context, account *models.Account) error {
	return r.DB.WithContext(ctx).Save(account).Error
}

func (r *accountRepository) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Delete(&models.Account{}, id).Error
}
