package http

import (
	"gobanking/internal/bank/dto"
	httpErrors "gobanking/pkg/http_errors"
	"gobanking/pkg/middlewares"
	"gobanking/pkg/payload"
	"gobanking/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TransactionHandler handler

func (h *TransactionHandler) Transfer(c echo.Context) error {
	res := &payload.Response{}
	req := &dto.TransferRequest{}
	if err := c.Bind(req); err != nil {
		h.Logger.Errorf("error binding request: %v", err)
		return httpErrors.NewBadRequestError(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	// validate request
	if err := h.validator.Validate(req); err != nil {
		h.Logger.Errorf("error validating request: %v", err)
		return httpErrors.NewBadRequestError(c, utils.GetErrorValidation(err), h.cfg.Http.DebugErrorsResponse)
	}

	// get user from context
	user := c.Get("user").(*middlewares.UserResponse)

	// transfer
	transaction, err := h.Usecase.Transaction.Transfer(c.Request().Context(), req, user.ID)
	if err != nil {
		h.Logger.Errorf("error transfer: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Data = transaction

	res.Message = "Success transfer"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}

func (h *TransactionHandler) Withdrawal(c echo.Context) error {
	res := &payload.Response{}
	req := &dto.WithdrawalRequest{}
	if err := c.Bind(req); err != nil {
		h.Logger.Errorf("error binding request: %v", err)
		return httpErrors.NewBadRequestError(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	// validate request
	if err := h.validator.Validate(req); err != nil {
		h.Logger.Errorf("error validating request: %v", err)
		return httpErrors.NewBadRequestError(c, utils.GetErrorValidation(err), h.cfg.Http.DebugErrorsResponse)
	}

	// get user from context
	user := c.Get("user").(*middlewares.UserResponse)

	// withdrawal
	transaction, err := h.Usecase.Transaction.Withdrawal(c.Request().Context(), req, user.ID)
	if err != nil {
		h.Logger.Errorf("error withdrawal: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Data = transaction

	res.Message = "Success withdrawal"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}

func (h *TransactionHandler) Deposit(c echo.Context) error {
	res := &payload.Response{}
	req := &dto.DepositRequest{}
	if err := c.Bind(req); err != nil {
		h.Logger.Errorf("error binding request: %v", err)
		return httpErrors.NewBadRequestError(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	// validate request
	if err := h.validator.Validate(req); err != nil {
		h.Logger.Errorf("error validating request: %v", err)
		return httpErrors.NewBadRequestError(c, utils.GetErrorValidation(err), h.cfg.Http.DebugErrorsResponse)
	}

	// get user from context
	user := c.Get("user").(*middlewares.UserResponse)

	// deposit
	transaction, err := h.Usecase.Transaction.Deposit(c.Request().Context(), req, user.ID)
	if err != nil {
		h.Logger.Errorf("error deposit: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Data = transaction

	res.Message = "Success deposit"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}
