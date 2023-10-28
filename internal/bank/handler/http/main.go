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
	Account  *AccountHandler
	Currency *CurrencyHandler
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
		Account: (*AccountHandler)(handler),
	}

	account := e.Group("/account")
	{
		account.POST("", h.Account.Create, handler.mw.AuthMiddlewareClient)
		account.GET("/:account_number", h.Account.GetByAccountNumber, handler.mw.AuthMiddlewareClient)
		account.GET("/user/:user_id", h.Account.GetByUserID, handler.mw.AuthMiddlewareClient)
		account.GET("/me", h.Account.GetByUserID, handler.mw.AuthMiddlewareClient)
		account.PUT("/:account_number", h.Account.Update, handler.mw.AuthMiddlewareClient)
		account.DELETE("/:account_number", h.Account.Delete, handler.mw.AuthMiddlewareClient)
	}

	currency := e.Group("/currency")
	{
		currency.GET("", h.Currency.GetAll)
		currency.GET("/:id", h.Currency.GetByID)
	}

	return h
}