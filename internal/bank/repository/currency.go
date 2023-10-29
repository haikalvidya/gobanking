package repository

import "gobanking/internal/bank/models"

type CurrencyRepository interface {
	// Get all currencies
	GetAll(limit, offset int) ([]*models.Currency, error)
	// Get currency by id
	GetByID(id int) (*models.Currency, error)
	// Get currency by code
	GetByCode(code string) (*models.Currency, error)
}

type currencyRepository repository

func (r *currencyRepository) GetAll(limit, offset int) ([]*models.Currency, error) {
	var currencies []*models.Currency
	err := r.DB.Limit(limit).Offset(offset).Find(&currencies).Error
	return currencies, err
}

func (r *currencyRepository) GetByID(id int) (*models.Currency, error) {
	currency := &models.Currency{}
	err := r.DB.Where("id = ?", id).First(currency).Error
	return currency, err
}

func (r *currencyRepository) GetByCode(code string) (*models.Currency, error) {
	currency := &models.Currency{}
	err := r.DB.Where("code = ?", code).First(currency).Error
	return currency, err
}
