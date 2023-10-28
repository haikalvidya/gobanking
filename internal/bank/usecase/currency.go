package usecase

import (
	"fmt"
	"gobanking/internal/bank/models"
	"gobanking/pkg/http_errors"
	"gobanking/pkg/payload"

	"gorm.io/gorm"
)

type CurrencyUsecase interface {
	// GetAll returns all currencies
	GetAll(pagination *payload.PaginationRequest) ([]*models.Currency, error)
	// GetByID returns a currency by id
	GetByID(id int) (*models.Currency, error)
}

type currencyUsecase usecase

func (u *currencyUsecase) GetAll(pagination *payload.PaginationRequest) ([]*models.Currency, error) {
	currencies, err := u.Repo.Currency.GetAll(int(pagination.Limit), int(pagination.Offset))
	if err != nil {
		u.Logger.Errorf("error when get all currencies: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when get all currencies: %v", typeOfErr, err)
	}

	return currencies, nil
}

func (u *currencyUsecase) GetByID(id int) (*models.Currency, error) {
	currency, err := u.Repo.Currency.GetByID(id)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Logger.Errorf("error when get currency by id: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when get currency by id: %v", typeOfErr, err)
	}

	if err == gorm.ErrRecordNotFound {
		u.Logger.Errorf("currency not found")
		typeOfErr := http_errors.NotFound
		return nil, fmt.Errorf("%v : currency not found", typeOfErr)
	}

	return currency, nil
}
