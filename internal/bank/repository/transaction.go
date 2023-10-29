package repository

import (
	"context"
	"gobanking/internal/bank/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	// CreateTX creates a new transaction
	CreateTX(ctx context.Context, tx *gorm.DB, transaction *models.Transaction) error
	// GetBySourceAccountNumber returns all transaction by source account number
	GetBySourceAccountNumber(ctx context.Context, accountNumber string, theType []string, limit, offset int) ([]*models.Transaction, error)
	// GetById returns transaction by id
	GetById(ctx context.Context, id string) (*models.Transaction, error)
}

type transactionRepository repository

func (r *transactionRepository) CreateTX(ctx context.Context, tx *gorm.DB, transaction *models.Transaction) error {
	return tx.WithContext(ctx).Create(transaction).Error
}

func (r *transactionRepository) GetBySourceAccountNumber(ctx context.Context, accountNumber string, theType []string, limit, offset int) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	err := r.DB.WithContext(ctx).Where("source_account_id = ? AND type IN ?", accountNumber, theType).Limit(limit).Offset(offset).Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) GetById(ctx context.Context, id string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&transaction).Error
	return &transaction, err
}
