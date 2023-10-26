package http

import (
	"gobanking/internal/user/dto"
	httpErrors "gobanking/pkg/http_errors"
	"gobanking/pkg/payload"
	"gobanking/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler handler

func (h *AuthHandler) SignUp(c echo.Context) error {
	res := &payload.Response{}
	req := &dto.SignUpRequest{}
	if err := c.Bind(req); err != nil {
		h.Logger.Errorf("error binding request: %v", err)
		return httpErrors.NewBadRequestError(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	// validate request
	if err := h.validator.Validate(req); err != nil {
		h.Logger.Errorf("error validating request: %v", err)
		return httpErrors.NewBadRequestError(c, utils.GetErrorValidation(err), h.cfg.Http.DebugErrorsResponse)
	}

	// create user
	user, err := h.Usecase.Auth.SignUp(c.Request().Context(), req)
	if err != nil {
		h.Logger.Errorf("error creating user: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Data = user

	res.Message = "Success signup"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) SignIn(c echo.Context) error {
	res := &payload.Response{}
	req := &dto.SignInRequest{}
	if err := c.Bind(req); err != nil {
		h.Logger.Errorf("error binding request: %v", err)
		return httpErrors.NewBadRequestError(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	// validate request
	if err := h.validator.Validate(req); err != nil {
		h.Logger.Errorf("error validating request: %v", err)
		return httpErrors.NewBadRequestError(c, utils.GetErrorValidation(err), h.cfg.Http.DebugErrorsResponse)
	}

	// check username or email, should be one that is not empty
	if req.Username == "" && req.Email == "" {
		h.Logger.Errorf("username or email should not be empty")
		return httpErrors.NewBadRequestError(c, "username or email should not be empty", h.cfg.Http.DebugErrorsResponse)
	}

	// sign in user
	user, err := h.Usecase.Auth.SignIn(c.Request().Context(), req)
	if err != nil {
		h.Logger.Errorf("error signing in user: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Data = user

	res.Message = "Success signin"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	res := &payload.Response{}

	// get refresh token from Authorization header
	refreshToken := c.Request().Header.Get("Authorization")
	if refreshToken == "" {
		h.Logger.Errorf("refresh token is empty")
		return httpErrors.NewUnauthorizedError(c, "Unauthorized", h.cfg.Http.DebugErrorsResponse)
	}

	// refresh token
	token, err := h.Usecase.Auth.RefreshToken(c.Request().Context(), refreshToken)
	if err != nil {
		h.Logger.Errorf("error refreshing token: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Data = token

	res.Message = "Success refresh token"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) SignOut(c echo.Context) error {
	res := &payload.Response{}

	// get token from Authorization header
	token := c.Request().Header.Get("Authorization")

	// delete token from usecase
	if err := h.Usecase.Auth.SignOut(c.Request().Context(), token); err != nil {
		h.Logger.Errorf("error signing out user: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Message = "Success signout"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Me(c echo.Context) error {
	res := &payload.Response{}

	// get token from Authorization header
	token := c.Request().Header.Get("Authorization")

	// get user from usecase
	user, err := h.Usecase.Auth.Me(c.Request().Context(), token)
	if err != nil {
		h.Logger.Errorf("error getting user: %v", err)
		return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
	}

	res.Data = user

	res.Message = "Success get user"
	res.Success = true

	return c.JSON(http.StatusOK, res)
}
