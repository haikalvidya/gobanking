package http

import (
	"gobanking/internal/user/dto"
	httpErrors "gobanking/pkg/http_errors"
	"gobanking/pkg/middlewares"
	"gobanking/pkg/payload"
	"gobanking/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler handler

func (h *UserHandler) GetByID(c echo.Context) error {
	res := &payload.Response{}

	// get user id from context
	userId := c.Get("user").(*middlewares.UserModel).ID

	// get user id from path
	id := c.Param("id")

	if userId.String() != id {
		h.Logger.Errorf("user id from context and path is not the same")
		return httpErrors.NewUnauthorizedError(c, "Unauthorized", h.cfg.Http.DebugErrorsResponse)
	}

	// get user from usecase
	user, err := h.Usecase.User.GetByID(c.Request().Context(), id)
	if err != nil {
		h.Logger.Errorf("error when get user by id: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Data = user

	res.Message = "Success get user by id"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}

func (h *UserHandler) Update(c echo.Context) error {
	res := &payload.Response{}
	req := &dto.UpdateUserRequest{}

	// bind request
	if err := c.Bind(req); err != nil {
		h.Logger.Errorf("error binding request: %v", err)
		return httpErrors.NewBadRequestError(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	// validate request
	if err := h.validator.Validate(req); err != nil {
		h.Logger.Errorf("error validating request: %v", err)
		return httpErrors.NewBadRequestError(c, utils.GetErrorValidation(err), h.cfg.Http.DebugErrorsResponse)
	}

	// get token from Authorization header
	token := c.Request().Header.Get("Authorization")

	// get user id from context
	userId := c.Get("user").(*middlewares.UserModel).ID

	// get user id from path
	id := c.Param("id")

	if userId.String() != id {
		h.Logger.Errorf("user id from context and path is not the same")
		return httpErrors.NewUnauthorizedError(c, "Unauthorized", h.cfg.Http.DebugErrorsResponse)
	}

	// get user from usecase
	user, err := h.Usecase.User.Update(c.Request().Context(), id, token, req)
	if err != nil {
		h.Logger.Errorf("error when update user: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Data = user

	res.Message = "Success update user"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}

func (h *UserHandler) Delete(c echo.Context) error {
	res := &payload.Response{}

	// get token from Authorization header
	token := c.Request().Header.Get("Authorization")

	// get user id from context
	userId := c.Get("user").(*middlewares.UserModel).ID

	// get user id from path
	id := c.Param("id")

	if userId.String() != id {
		h.Logger.Errorf("user id from context and path is not the same")
		return httpErrors.NewUnauthorizedError(c, "Unauthorized", h.cfg.Http.DebugErrorsResponse)
	}

	// delete user from usecase
	if err := h.Usecase.User.Delete(c.Request().Context(), id, token); err != nil {
		h.Logger.Errorf("error when delete user: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Message = "Success delete user"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}
