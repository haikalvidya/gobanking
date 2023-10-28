package http

import (
	httpErrors "gobanking/pkg/http_errors"
	"gobanking/pkg/payload"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CurrencyHandler handler

func (h *CurrencyHandler) GetAll(c echo.Context) error {
	res := &payload.Response{}

	// get pagination
	// page and per_page query
	pagination := payload.GetPagination(c)

	// get all currencies
	currencies, err := h.Usecase.Currency.GetAll(pagination)
	if err != nil {
		h.Logger.Errorf("error when get all currencies: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Data = currencies

	res.Message = "Success get all currencies"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}

func (h *CurrencyHandler) GetByID(c echo.Context) error {
	res := &payload.Response{}

	// get currency id from path
	id := c.Param("id")
	// convert id to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Errorf("error when convert id to int: %v", err)
		return httpErrors.NewBadRequestError(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	// get currency from usecase
	currency, err := h.Usecase.Currency.GetByID(idInt)
	if err != nil {
		h.Logger.Errorf("error when get currency by id: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Data = currency

	res.Message = "Success get currency by id"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}
