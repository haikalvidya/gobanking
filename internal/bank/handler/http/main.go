package http

import (
	"gobanking/internal/bank/config"
	"gobanking/internal/bank/usecase"
	"gobanking/pkg/logger"
	"gobanking/pkg/middlewares"
	"gobanking/pkg/utils"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Account     *AccountHandler
	Currency    *CurrencyHandler
	Transaction *TransactionHandler
}

type handler struct {
	Usecase   *usecase.Usecase
	Logger    logger.Logger
	cfg       *config.Config
	mw        middlewares.MiddlewareManager
	validator *utils.CustomValidator
}

func NewHandler(usecase *usecase.Usecase,
	logger logger.Logger,
	cfg *config.Config,
	mw middlewares.MiddlewareManager,
	validator *utils.CustomValidator,
	e *echo.Group,
) *Handler {
	handler := &handler{
		Usecase:   usecase,
		Logger:    logger,
		cfg:       cfg,
		mw:        mw,
		validator: validator,
	}

	h := &Handler{
		Account:     (*AccountHandler)(handler),
		Currency:    (*CurrencyHandler)(handler),
		Transaction: (*TransactionHandler)(handler),
	}

	account := e.Group("/account")
	{
		account.POST("", h.Account.Create, handler.mw.AuthMiddlewareClient)
		account.GET("/me", h.Account.GetByME, handler.mw.AuthMiddlewareClient)
		account.GET("/detail/:account_number", h.Account.GetByAccountNumberByOwner, handler.mw.AuthMiddlewareClient)
		account.PUT("/:account_number", h.Account.Update, handler.mw.AuthMiddlewareClient)
		account.DELETE("/:account_number", h.Account.Delete, handler.mw.AuthMiddlewareClient)
		account.GET("/:account_number", h.Account.GetByAccountNumberExternal)
	}

	currency := e.Group("/currency")
	{
		currency.GET("", h.Currency.GetAll)
		currency.GET("/:id", h.Currency.GetByID)
	}

	transaction := e.Group("/transaction")
	{
		transaction.POST("/transfer", h.Transaction.Transfer, handler.mw.AuthMiddlewareClient)
		transaction.POST("/withdrawal", h.Transaction.Withdrawal, handler.mw.AuthMiddlewareClient)
		transaction.POST("/deposit", h.Transaction.Deposit, handler.mw.AuthMiddlewareClient)
	}

	return h
}
