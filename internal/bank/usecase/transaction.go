package usecase

import (
	"context"
	"fmt"
	"gobanking/internal/bank/dto"
	"gobanking/internal/bank/models"
	"gobanking/pkg/http_errors"
	"gobanking/pkg/payload"

	"github.com/Rhymond/go-money"
	"gorm.io/gorm"
)

type TransactionUsecase interface {
	// Transfer transfers money from source account to destination account
	Transfer(ctx context.Context, transferReq *dto.TransferRequest, userId string) (*models.Transaction, error)
	// Withdrawal withdraws money from source account
	Withdrawal(ctx context.Context, withdrawalReq *dto.WithdrawalRequest, userId string) (*models.Transaction, error)
	// Deposit deposits money to destination account
	Deposit(ctx context.Context, depositReq *dto.DepositRequest, userId string) (*models.Transaction, error)
	// GetBySourceAccountNumber returns all transaction by source account number
	GetBySourceAccountNumber(ctx context.Context, accountNumber string, theType []string, pagination *payload.PaginationRequest) ([]*models.Transaction, error)
}

type transactionUsecase usecase

func (u *transactionUsecase) Transfer(ctx context.Context, transferReq *dto.TransferRequest, userId string) (*models.Transaction, error) {
	// check if the amount is valid
	if transferReq.Amount <= 0 {
		u.Logger.Errorf("amount must be greater than 0")
		typeOfErr := http_errors.BadRequest
		return nil, fmt.Errorf("%v : amount must be greater than 0", typeOfErr)
	}

	// check if source account number is valid
	sourceAccount, err := u.Repo.Account.GetByAccountNumber(ctx, transferReq.SourceAccountNumber)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Logger.Errorf("error when get source account by account number: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when get source account by account number: %v", typeOfErr, err)
	}

	if err == gorm.ErrRecordNotFound {
		u.Logger.Errorf("source account not found")
		typeOfErr := http_errors.NotFound
		return nil, fmt.Errorf("%v : source account not found", typeOfErr)
	}

	// check if source account belongs to user
	if sourceAccount.UserId.String() != userId {
		u.Logger.Errorf("source account does not belong to user")
		typeOfErr := http_errors.Forbidden
		return nil, fmt.Errorf("%v : source account does not belong to user", typeOfErr)
	}

	// check if destination account number is valid
	destinationAccount, err := u.Repo.Account.GetByAccountNumber(ctx, transferReq.DestinationAccountNumber)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Logger.Errorf("error when get destination account by account number: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when get destination account by account number: %v", typeOfErr, err)
	}

	if err == gorm.ErrRecordNotFound {
		u.Logger.Errorf("destination account not found")
		typeOfErr := http_errors.NotFound
		return nil, fmt.Errorf("%v : destination account not found", typeOfErr)
	}

	// check if source account and destination account are the same
	if sourceAccount.AccountNumber == destinationAccount.AccountNumber {
		u.Logger.Errorf("source account and destination account are the same")
		typeOfErr := http_errors.BadRequest
		return nil, fmt.Errorf("%v : source account and destination account are the same", typeOfErr)
	}

	// check if source account and destination account are the same currency
	if sourceAccount.CurrencyId != destinationAccount.CurrencyId {
		u.Logger.Errorf("source account and destination account are not the same currency")
		typeOfErr := http_errors.BadRequest
		return nil, fmt.Errorf("%v : source account and destination account are not the same currency", typeOfErr)
	}

	// get code from currency id
	theCurrencyModel, err := u.Repo.Currency.GetByID(sourceAccount.CurrencyId)
	if err != nil {
		u.Logger.Errorf("error when get currency by id: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when get currency by id: %v", typeOfErr, err)
	}

	// currency model
	currency := money.GetCurrency(theCurrencyModel.Code)
	if currency == nil {
		u.Logger.Errorf("currency not found")
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : currency not found", typeOfErr)
	}

	transactionId := ""

	// do in transaction
	txFunc := func(tx *gorm.DB) error {
		// convert money-related value int64 to money.Money
		sourceAccountBalance := money.New(sourceAccount.Balance, currency.Code)
		transferAmount := money.New(transferReq.Amount, currency.Code)
		destAccountBalance := money.New(destinationAccount.Balance, currency.Code)

		isEnoughBalance, err := sourceAccountBalance.GreaterThanOrEqual(transferAmount)
		if err != nil {
			u.Logger.Errorf("error when check if source account has enough balance: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when check if source account has enough balance: %v", typeOfErr, err)
		}
		// check if source account has enough balance
		if !isEnoughBalance {
			u.Logger.Errorf("source account has not enough balance")
			typeOfErr := http_errors.BadRequest
			return fmt.Errorf("%v : source account has not enough balance", typeOfErr)
		}

		// update source account balance
		sourceAccountBalance, err = sourceAccountBalance.Subtract(transferAmount)
		if err != nil {
			u.Logger.Errorf("error when subtract source account balance: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when subtract source account balance: %v", typeOfErr, err)
		}
		sourceAccount.Balance = sourceAccountBalance.Amount()
		err = u.Repo.Account.UpdateTX(ctx, tx, sourceAccount)
		if err != nil {
			u.Logger.Errorf("error when update source account balance: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when update source account balance: %v", typeOfErr, err)
		}

		// update destination account balance
		destAccountBalance, err = destAccountBalance.Add(transferAmount)
		if err != nil {
			u.Logger.Errorf("error when add destination account balance: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when add destination account balance: %v", typeOfErr, err)
		}
		destinationAccount.Balance = destAccountBalance.Amount()
		err = u.Repo.Account.UpdateTX(ctx, tx, destinationAccount)
		if err != nil {
			u.Logger.Errorf("error when update destination account balance: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when update destination account balance: %v", typeOfErr, err)
		}

		// create transaction
		transaction := &models.Transaction{
			SourceAccountNumber: sourceAccount.AccountNumber,
			DestAccountNumber:   destinationAccount.AccountNumber,
			Type:                models.TRANSFER,
			Amount:              transferReq.Amount,
		}
		err = u.Repo.Transaction.CreateTX(ctx, tx, transaction)
		if err != nil {
			u.Logger.Errorf("error when create transaction: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when create transaction: %v", typeOfErr, err)
		}

		transactionId = transaction.ID.String()

		// save the ledger
		ledgerSourceAccount := &models.BalanceLedger{
			AccountNumber: sourceAccount.AccountNumber,
			Balance:       sourceAccount.Balance,
			TransactionId: transaction.ID,
		}
		err = u.Repo.BalanceLedger.CreateTX(ctx, tx, ledgerSourceAccount)
		if err != nil {
			u.Logger.Errorf("error when create ledger for source account: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when create ledger for source account: %v", typeOfErr, err)
		}

		ledgerDestAccount := &models.BalanceLedger{
			AccountNumber: destinationAccount.AccountNumber,
			Balance:       destinationAccount.Balance,
			TransactionId: transaction.ID,
		}
		err = u.Repo.BalanceLedger.CreateTX(ctx, tx, ledgerDestAccount)
		if err != nil {
			u.Logger.Errorf("error when create ledger for destination account: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when create ledger for destination account: %v", typeOfErr, err)
		}

		return nil
	}

	err = u.Repo.Tx.DoInTransaction(txFunc)
	if err != nil {
		return nil, err
	}

	resp, err := u.Repo.Transaction.GetById(ctx, transactionId)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Logger.Errorf("error when get transaction by id: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when get transaction by id: %v", typeOfErr, err)
	}

	return resp, nil
}

func (u *transactionUsecase) Withdrawal(ctx context.Context, withdrawalReq *dto.WithdrawalRequest, userId string) (*models.Transaction, error) {
	// check if the amount is valid
	if withdrawalReq.Amount <= 0 {
		u.Logger.Errorf("amount must be greater than 0")
		typeOfErr := http_errors.BadRequest
		return nil, fmt.Errorf("%v : amount must be greater than 0", typeOfErr)
	}

	// check if source account number is valid
	sourceAccount, err := u.checkIfAccountBelongsToUser(ctx, withdrawalReq.SourceAccountNumber, userId)
	if err != nil {
		return nil, err
	}

	// get code from currency id
	theCurrencyModel, err := u.Repo.Currency.GetByID(sourceAccount.CurrencyId)
	if err != nil {
		u.Logger.Errorf("error when get currency by id: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when get currency by id: %v", typeOfErr, err)
	}

	// currency model
	currency := money.GetCurrency(theCurrencyModel.Code)
	if currency == nil {
		u.Logger.Errorf("currency not found")
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : currency not found", typeOfErr)
	}

	transactionId := ""

	// do in transaction
	txFunc := func(tx *gorm.DB) error {
		// convert money-related value int64 to money.Money
		sourceAccountBalance := money.New(sourceAccount.Balance, currency.Code)
		withdrawalAmount := money.New(withdrawalReq.Amount, currency.Code)

		isEnoughBalance, err := sourceAccountBalance.GreaterThanOrEqual(withdrawalAmount)
		if err != nil {
			u.Logger.Errorf("error when check if source account has enough balance: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when check if source account has enough balance: %v", typeOfErr, err)
		}
		// check if source account has enough balance
		if !isEnoughBalance {
			u.Logger.Errorf("source account has not enough balance")
			typeOfErr := http_errors.BadRequest
			return fmt.Errorf("%v : source account has not enough balance", typeOfErr)
		}

		// update source account balance
		sourceAccountBalance, err = sourceAccountBalance.Subtract(withdrawalAmount)
		if err != nil {
			u.Logger.Errorf("error when subtract source account balance: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when subtract source account balance: %v", typeOfErr, err)
		}
		sourceAccount.Balance = sourceAccountBalance.Amount()
		err = u.Repo.Account.UpdateTX(ctx, tx, sourceAccount)
		if err != nil {
			u.Logger.Errorf("error when update source account balance: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when update source account balance: %v", typeOfErr, err)
		}

		// create transaction
		transaction := &models.Transaction{
			SourceAccountNumber: sourceAccount.AccountNumber,
			Type:                models.WITHDRAWAL,
			Amount:              withdrawalReq.Amount,
		}
		err = u.Repo.Transaction.CreateTX(ctx, tx, transaction)
		if err != nil {
			u.Logger.Errorf("error when create transaction: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when create transaction: %v", typeOfErr, err)
		}

		transactionId = transaction.ID.String()

		// save the ledger
		ledgerSourceAccount := &models.BalanceLedger{
			AccountNumber: sourceAccount.AccountNumber,
			Balance:       sourceAccount.Balance,
			TransactionId: transaction.ID,
		}
		err = u.Repo.BalanceLedger.CreateTX(ctx, tx, ledgerSourceAccount)
		if err != nil {
			u.Logger.Errorf("error when create ledger for source account: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when create ledger for source account: %v", typeOfErr, err)
		}

		return nil
	}

	err = u.Repo.Tx.DoInTransaction(txFunc)
	if err != nil {
		return nil, err
	}

	resp, err := u.Repo.Transaction.GetById(ctx, transactionId)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Logger.Errorf("error when get transaction by id: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when get transaction by id: %v", typeOfErr, err)
	}

	return resp, nil
}

func (u *transactionUsecase) Deposit(ctx context.Context, depositReq *dto.DepositRequest, userId string) (*models.Transaction, error) {
	// check if the amount is valid
	if depositReq.Amount <= 0 {
		u.Logger.Errorf("amount must be greater than 0")
		typeOfErr := http_errors.BadRequest
		return nil, fmt.Errorf("%v : amount must be greater than 0", typeOfErr)
	}

	// check if destination account number is valid
	destAccount, err := u.checkIfAccountBelongsToUser(ctx, depositReq.DestinationAccountNumber, userId)
	if err != nil {
		return nil, err
	}

	// get code from currency id
	theCurrencyModel, err := u.Repo.Currency.GetByID(destAccount.CurrencyId)
	if err != nil {
		u.Logger.Errorf("error when get currency by id: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when get currency by id: %v", typeOfErr, err)
	}

	// currency model
	currency := money.GetCurrency(theCurrencyModel.Code)
	if currency == nil {
		u.Logger.Errorf("currency not found")
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : currency not found", typeOfErr)
	}

	transactionId := ""

	// do in transaction
	txFunc := func(tx *gorm.DB) error {
		// convert money-related value int64 to money.Money
		destAccountBalance := money.New(destAccount.Balance, currency.Code)
		depositAmount := money.New(depositReq.Amount, currency.Code)

		// update destination account balance
		destAccountBalance, err = destAccountBalance.Add(depositAmount)
		if err != nil {
			u.Logger.Errorf("error when add destination account balance: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when add destination account balance: %v", typeOfErr, err)
		}
		destAccount.Balance = destAccountBalance.Amount()
		err = u.Repo.Account.UpdateTX(ctx, tx, destAccount)
		if err != nil {
			u.Logger.Errorf("error when update destination account balance: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when update destination account balance: %v", typeOfErr, err)
		}

		// create transaction
		transaction := &models.Transaction{
			DestAccountNumber: destAccount.AccountNumber,
			Type:              models.DEPOSIT,
			Amount:            depositReq.Amount,
		}
		err = u.Repo.Transaction.CreateTX(ctx, tx, transaction)
		if err != nil {
			u.Logger.Errorf("error when create transaction: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when create transaction: %v", typeOfErr, err)
		}

		transactionId = transaction.ID.String()

		// save the ledger
		ledgerDestAccount := &models.BalanceLedger{
			AccountNumber: destAccount.AccountNumber,
			Balance:       destAccount.Balance,
			TransactionId: transaction.ID,
		}
		err = u.Repo.BalanceLedger.CreateTX(ctx, tx, ledgerDestAccount)
		if err != nil {
			u.Logger.Errorf("error when create ledger for destination account: %v", err)
			typeOfErr := http_errors.InternalServerError
			return fmt.Errorf("%v : error when create ledger for destination account: %v", typeOfErr, err)
		}

		return nil
	}

	err = u.Repo.Tx.DoInTransaction(txFunc)
	if err != nil {
		return nil, err
	}

	resp, err := u.Repo.Transaction.GetById(ctx, transactionId)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Logger.Errorf("error when get transaction by id: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when get transaction by id: %v", typeOfErr, err)
	}

	return resp, nil
}

func (u *transactionUsecase) GetBySourceAccountNumber(ctx context.Context, accountNumber string, theType []string, pagination *payload.PaginationRequest) ([]*models.Transaction, error) {
	return nil, nil
}

func (u *transactionUsecase) checkIfAccountBelongsToUser(ctx context.Context, accountNumber string, userId string) (*models.Account, error) {
	// check if account number is valid
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

	// check if account belongs to user
	if account.UserId.String() != userId {
		u.Logger.Errorf("account does not belong to user")
		typeOfErr := http_errors.Forbidden
		return nil, fmt.Errorf("%v : account does not belong to user", typeOfErr)
	}

	return account, nil
}
