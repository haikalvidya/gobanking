package usecase

import (
	"context"
	"fmt"
	"gobanking/internal/bank/dto"
	"gobanking/internal/bank/models"
	"gobanking/pkg/http_errors"
	"gobanking/pkg/payload"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountUsecase interface {
	// Create creates a new account
	Create(ctx context.Context, accountReq *dto.CreateAccountRequest, userId string) (*models.Account, error)
	// GetByAccountNumberByOwner returns an account by account number and owner
	GetByAccountNumberByOwner(ctx context.Context, accountNumber string, userId string) (*models.Account, error)
	// GetByAccountNumberExternal returns an account by account number
	GetByAccountNumberExternal(ctx context.Context, accountNumber string) (*dto.GetAccountByExternalResponse, error)
	// GetByUserID returns all account by user id
	GetByUserID(ctx context.Context, userID string, pagination *payload.PaginationRequest) ([]*models.Account, error)
	// Update updates an account
	Update(ctx context.Context, accountNumber string, account *dto.UpdateAccountRequest, userId string) (*models.Account, error)
	// Delete deletes an account
	Delete(ctx context.Context, accountNumber string, userId string) error
}

type accountUsecase usecase

func (u *accountUsecase) Create(ctx context.Context, accountReq *dto.CreateAccountRequest, userId string) (*models.Account, error) {
	resp := &models.Account{}
	err := u.Repo.Tx.DoInTransaction(func(tx *gorm.DB) error {
		// check if the currency is valid
		currency, err := u.Repo.Currency.GetByCode(accountReq.CurrencyCode)
		if err != nil && err != gorm.ErrRecordNotFound {
			u.Logger.Errorf("error when get currency by code: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when get currency by code: %v", typeOfErr, err)
		}

		if err == gorm.ErrRecordNotFound {
			u.Logger.Errorf("currency not found")
			typeOfErr := http_errors.NotFound
			return fmt.Errorf("%v : currency not found", typeOfErr)
		}

		// convert user id to uuid
		userIdUUID, err := uuid.Parse(userId)
		if err != nil {
			u.Logger.Errorf("error when parse user id to uuid: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when parse user id to uuid: %v", typeOfErr, err)
		}

		// create account
		account := &models.Account{
			UserId:     userIdUUID,
			Name:       accountReq.Name,
			CurrencyId: currency.ID,
		}

		err = u.Repo.Account.CreateTX(ctx, tx, account)
		if err != nil {
			u.Logger.Errorf("error when create account: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when create account: %v", typeOfErr, err)
		}

		resp = account

		return nil
	})

	return resp, err
}

func (u *accountUsecase) GetByAccountNumberExternal(ctx context.Context, accountNumber string) (*dto.GetAccountByExternalResponse, error) {
	account, err := u.Repo.Account.GetByAccountNumber(ctx, accountNumber)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Logger.Errorf("error when get account by account number: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when get account by account number: %v", typeOfErr, err)
	}

	if err == gorm.ErrRecordNotFound {
		u.Logger.Errorf("account not found")
		typeOfErr := http_errors.NotFound
		return nil, fmt.Errorf("%v : account not found", typeOfErr)
	}

	resp := &dto.GetAccountByExternalResponse{
		AccountNumber: account.AccountNumber,
		Name:          account.Name,
		CurrencyId:    account.CurrencyId,
		CreatedAt:     account.CreatedAt,
		UpdatedAt:     account.UpdatedAt,
	}

	return resp, nil
}

func (u *accountUsecase) GetByAccountNumberByOwner(ctx context.Context, accountNumber string, userId string) (*models.Account, error) {
	// check if the account is owned by the user
	account, err := u.Repo.Account.GetByAccountNumberAndUserID(ctx, accountNumber, userId)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Logger.Errorf("error when get account by account number and user id: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when get account by account number and user id: %v", typeOfErr, err)
	}

	if err == gorm.ErrRecordNotFound {
		u.Logger.Errorf("account not found")
		typeOfErr := http_errors.NotFound
		return nil, fmt.Errorf("%v : account not found", typeOfErr)
	}

	return account, nil
}

func (u *accountUsecase) GetByUserID(ctx context.Context, userID string, pagination *payload.PaginationRequest) ([]*models.Account, error) {
	limit := int(pagination.Limit)
	offset := int(pagination.Offset)
	accounts, err := u.Repo.Account.GetByUserID(ctx, userID, limit, offset)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Logger.Errorf("error when get account by user id: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when get account by user id: %v", typeOfErr, err)
	}

	if err == gorm.ErrRecordNotFound {
		u.Logger.Errorf("account not found")
		typeOfErr := http_errors.NotFound
		return nil, fmt.Errorf("%v : account not found", typeOfErr)
	}

	return accounts, nil
}

func (u *accountUsecase) Update(ctx context.Context, accountNumber string, account *dto.UpdateAccountRequest, userId string) (*models.Account, error) {
	resp := &models.Account{}
	err := u.Repo.Tx.DoInTransaction(func(tx *gorm.DB) error {
		// get the account by account number
		accountModel, err := u.Repo.Account.GetByAccountNumber(ctx, accountNumber)
		if err != nil && err != gorm.ErrRecordNotFound {
			u.Logger.Errorf("error when get account by account number: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when get account by account number: %v", typeOfErr, err)
		}

		if err == gorm.ErrRecordNotFound {
			u.Logger.Errorf("account not found")
			typeOfErr := http_errors.NotFound
			return fmt.Errorf("%v : account not found", typeOfErr)
		}

		// check if the account is owned by the user
		if accountModel.UserId.String() != userId {
			u.Logger.Errorf("account not found")
			typeOfErr := http_errors.NotFound
			return fmt.Errorf("%v : account not found", typeOfErr)
		}

		// update account
		if account.Name != "" {
			accountModel.Name = account.Name
		}

		err = u.Repo.Account.UpdateTX(ctx, tx, accountModel)
		if err != nil {
			u.Logger.Errorf("error when update account: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when update account: %v", typeOfErr, err)
		}

		if err == gorm.ErrRecordNotFound {
			u.Logger.Errorf("account not found")
			typeOfErr := http_errors.NotFound
			return fmt.Errorf("%v : account not found", typeOfErr)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// get the account by account number
	resp, err = u.Repo.Account.GetByAccountNumber(ctx, accountNumber)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Logger.Errorf("error when get account by account number: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when get account by account number: %v", typeOfErr, err)
	}

	if err == gorm.ErrRecordNotFound {
		u.Logger.Errorf("account not found")
		typeOfErr := http_errors.NotFound
		return nil, fmt.Errorf("%v : account not found", typeOfErr)
	}

	return resp, nil
}

func (u *accountUsecase) Delete(ctx context.Context, accountNumber string, userId string) error {
	// check if the account is owned by the user
	_, err := u.Repo.Account.GetByAccountNumberAndUserID(ctx, accountNumber, userId)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Logger.Errorf("error when get account by account number and user id: %v", err)
		typeOfErr := http_errors.InternalServerError
		return fmt.Errorf("%v : error when get account by account number and user id: %v", typeOfErr, err)
	}

	if err == gorm.ErrRecordNotFound {
		u.Logger.Errorf("account not found")
		typeOfErr := http_errors.NotFound
		return fmt.Errorf("%v : account not found", typeOfErr)
	}

	err = u.Repo.Account.Delete(ctx, accountNumber)
	if err != nil {
		u.Logger.Errorf("error when delete account: %v", err)
		typeOfErr := http_errors.InternalServerError
		return fmt.Errorf("%v : error when delete account: %v", typeOfErr, err)
	}

	return err
}
