package http

import (
	"gobanking/internal/bank/dto"
	httpErrors "gobanking/pkg/http_errors"
	"gobanking/pkg/payload"
	"gobanking/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AccountHandler handler

func (h *AccountHandler) Create(c echo.Context) error {
	res := &payload.Response{}
	req := &dto.CreateAccountRequest{}
	if err := c.Bind(req); err != nil {
		h.Logger.Errorf("error binding request: %v", err)
		return httpErrors.NewBadRequestError(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	// validate request
	if err := h.validator.Validate(req); err != nil {
		h.Logger.Errorf("error validating request: %v", err)
		return httpErrors.NewBadRequestError(c, utils.GetErrorValidation(err), h.cfg.Http.DebugErrorsResponse)
	}

	// create account
	account, err := h.Usecase.Account.Create(c.Request().Context(), req)
	if err != nil {
		h.Logger.Errorf("error creating account: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Data = account

	res.Message = "Success create account"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}

func (h *AccountHandler) GetByAccountNumber(c echo.Context) error {
	res := &payload.Response{}

	// get account number from path
	accountNumber := c.Param("account_number")

	// get account from usecase
	account, err := h.Usecase.Account.GetByAccountNumber(c.Request().Context(), accountNumber)
	if err != nil {
		h.Logger.Errorf("error when get account by account number: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Data = account

	res.Message = "Success get account by account number"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}

func (h *AccountHandler) GetByUserID(c echo.Context) error {
	res := &payload.Response{}

	// get user id from path
	userID := c.Param("user_id")

	// get account from usecase
	account, err := h.Usecase.Account.GetByUserID(c.Request().Context(), userID)
	if err != nil {
		h.Logger.Errorf("error when get account by user id: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Data = account

	res.Message = "Success get account by user id"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}

func (h *AccountHandler) Update(c echo.Context) error {
	res := &payload.Response{}
	req := &dto.UpdateAccountRequest{}
	if err := c.Bind(req); err != nil {
		h.Logger.Errorf("error binding request: %v", err)
		return httpErrors.NewBadRequestError(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	// validate request
	if err := h.validator.Validate(req); err != nil {
		h.Logger.Errorf("error validating request: %v", err)
		return httpErrors.NewBadRequestError(c, utils.GetErrorValidation(err), h.cfg.Http.DebugErrorsResponse)
	}

	// get account number from path
	accountNumber := c.Param("account_number")

	// update account
	account, err := h.Usecase.Account.Update(c.Request().Context(), accountNumber, req)
	if err != nil {
		h.Logger.Errorf("error updating account: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Data = account

	res.Message = "Success update account"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}

func (h *AccountHandler) Delete(c echo.Context) error {
	res := &payload.Response{}

	// get id from path
	accountNumber := c.Param("account_number")

	// delete account
	err := h.Usecase.Account.Delete(c.Request().Context(), accountNumber)
	if err != nil {
		h.Logger.Errorf("error deleting account: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Message = "Success delete account"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}
