package usecase

import (
	"context"
	"gobanking/internal/bank/dto"
	"gobanking/internal/bank/models"
)

type AccountUsecase interface {
	// Create creates a new account
	Create(ctx context.Context, accountReq *dto.CreateAccountRequest) (*models.Account, error)
	// GetByAccountNumber returns an account by account number
	GetByAccountNumber(ctx context.Context, accountNumber string) (*models.Account, error)
	// GetByUserID returns all account by user id
	GetByUserID(ctx context.Context, userID string) ([]*models.Account, error)
	// Update updates an account
	Update(ctx context.Context, accountNumber string, account *dto.UpdateAccountRequest) (*models.Account, error)
	// Delete deletes an account
	Delete(ctx context.Context, id string) error
}

type accountUsecase usecase

func (u *accountUsecase) Create(ctx context.Context, accountReq *dto.CreateAccountRequest) (*models.Account, error) {
	return nil, nil
}

func (u *accountUsecase) GetByAccountNumber(ctx context.Context, accountNumber string) (*models.Account, error) {
	return nil, nil
}

func (u *accountUsecase) GetByUserID(ctx context.Context, userID string) ([]*models.Account, error) {
	return nil, nil
}

func (u *accountUsecase) Update(ctx context.Context, accountNumber string, account *dto.UpdateAccountRequest) (*models.Account, error) {
	return nil, nil
}

func (u *accountUsecase) Delete(ctx context.Context, id string) error {
	return nil
}
