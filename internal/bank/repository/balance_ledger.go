package repository

import (
	"context"
	"gobanking/internal/bank/models"

	"gorm.io/gorm"
)

type BalanceLedgerRepository interface {
	// CreateTX creates a new balance ledger
	CreateTX(ctx context.Context, tx *gorm.DB, balanceLedger *models.BalanceLedger) error
	// GetByAccountNumber returns all balance ledger by account number
	GetByAccountNumber(ctx context.Context, accountNumber string, limit, offset int) ([]*models.BalanceLedger, error)
}

type balanceLedgerRepository repository

func (r *balanceLedgerRepository) CreateTX(ctx context.Context, tx *gorm.DB, balanceLedger *models.BalanceLedger) error {
	return tx.WithContext(ctx).Create(balanceLedger).Error
}

func (r *balanceLedgerRepository) GetByAccountNumber(ctx context.Context, accountNumber string, limit, offset int) ([]*models.BalanceLedger, error) {
	var balanceLedgers []*models.BalanceLedger
	err := r.DB.WithContext(ctx).Where("account_number = ?", accountNumber).Limit(limit).Offset(offset).Find(&balanceLedgers).Error
	return balanceLedgers, err
}
