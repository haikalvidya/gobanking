package usecase

import (
	"gobanking/internal/bank/repository"
	"gobanking/pkg/logger"
)

type Usecase struct {
	Account     AccountUsecase
	Currency    CurrencyUsecase
	Transaction TransactionUsecase
}

type usecase struct {
	Repo   *repository.Repository
	Logger logger.Logger
}

func NewUsecase(repo *repository.Repository, logger logger.Logger) *Usecase {
	svc := &usecase{
		Repo:   repo,
		Logger: logger,
	}
	return &Usecase{
		Account:     (*accountUsecase)(svc),
		Currency:    (*currencyUsecase)(svc),
		Transaction: (*transactionUsecase)(svc),
	}
}
