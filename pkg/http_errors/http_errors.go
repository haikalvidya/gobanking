package http_errors

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

const (
	ErrBadRequest          = "Bad request"
	ErrNotFound            = "Not Found"
	ErrUnauthorized        = "Unauthorized"
	ErrRequestTimeout      = "Request Timeout"
	ErrInvalidEmail        = "Invalid email"
	ErrInvalidPassword     = "Invalid password"
	ErrInvalidField        = "Invalid field"
	ErrInternalServerError = "Internal Server Error"
	ErrWrongCredentials    = "Wrong Credentials"
	ErrForbidden           = "Forbidden"
)

var (
	BadRequest          = errors.New(ErrBadRequest)
	WrongCredentials    = errors.New(ErrWrongCredentials)
	NotFound            = errors.New(ErrNotFound)
	Unauthorized        = errors.New(ErrUnauthorized)
	Forbidden           = errors.New(ErrForbidden)
	InternalServerError = errors.New(ErrInternalServerError)
)

// RestErr Rest error interface
type RestErr interface {
	Status() int
	Error() string
	Causes() interface{}
	ErrBody() RestError
}

// RestError Rest error struct
type RestError struct {
	Success    bool        `json:"success"`
	ErrStatus  int         `json:"status,omitempty"`
	ErrError   string      `json:"error,omitempty"`
	ErrMessage interface{} `json:"message,omitempty"`
	Timestamp  time.Time   `json:"timestamp,omitempty"`
}

// ErrBody Error body
func (e RestError) ErrBody() RestError {
	return e
}

// Error  Error() interface method
func (e RestError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError, e.ErrMessage)
}

// Status Error status
func (e RestError) Status() int {
	return e.ErrStatus
}

// Causes RestError Causes
func (e RestError) Causes() interface{} {
	return e.ErrMessage
}

// Success RestError Success
func (e RestError) IsSuccess() bool {
	return e.Success
}

// NewRestError New Rest Error
func NewRestError(status int, err string, causes interface{}, debug bool) RestErr {
	restError := RestError{
		Success:   false,
		ErrStatus: status,
		ErrError:  err,
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	return restError
}

// NewRestErrorWithMessage New Rest Error With Message
func NewRestErrorWithMessage(status int, err string, causes interface{}) RestErr {
	return RestError{
		Success:    false,
		ErrStatus:  status,
		ErrError:   err,
		ErrMessage: causes,
		Timestamp:  time.Now().UTC(),
	}
}

// NewBadRequestError New Bad Request Error
func NewBadRequestError(ctx echo.Context, causes interface{}, debug bool) error {
	restError := RestError{
		Success:   false,
		ErrStatus: http.StatusBadRequest,
		ErrError:  BadRequest.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	return ctx.JSON(http.StatusBadRequest, restError)
}

// NewNotFoundError New Not Found Error
func NewNotFoundError(ctx echo.Context, causes interface{}, debug bool) error {
	restError := RestError{
		Success:   false,
		ErrStatus: http.StatusNotFound,
		ErrError:  NotFound.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	return ctx.JSON(http.StatusNotFound, restError)
}

// NewUnauthorizedError New Unauthorized Error
func NewUnauthorizedError(ctx echo.Context, causes interface{}, debug bool) error {

	restError := RestError{
		Success:   false,
		ErrStatus: http.StatusUnauthorized,
		ErrError:  Unauthorized.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	return ctx.JSON(http.StatusUnauthorized, restError)
}

// NewMethodNotAllowedError New Method Not Allowed Error
func NewMethodNotAllowedError(ctx echo.Context, causes interface{}, debug bool) error {

	restError := RestError{
		Success:   false,
		ErrStatus: http.StatusMethodNotAllowed,
		ErrError:  http.StatusText(http.StatusMethodNotAllowed),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	return ctx.JSON(http.StatusMethodNotAllowed, restError)
}

// NewForbiddenError New Forbidden Error
func NewForbiddenError(ctx echo.Context, causes interface{}, debug bool) error {

	restError := RestError{
		Success:   false,
		ErrStatus: http.StatusForbidden,
		ErrError:  Forbidden.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	return ctx.JSON(http.StatusForbidden, restError)
}

// NewInternalServerError New Internal Server Error
func NewInternalServerError(ctx echo.Context, causes interface{}, debug bool) error {

	restError := RestError{
		Success:   false,
		ErrStatus: http.StatusInternalServerError,
		ErrError:  InternalServerError.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	return ctx.JSON(http.StatusInternalServerError, restError)
}

// ParseErrors Parser of error string messages returns RestError
func ParseErrors(err error, debug bool) RestErr {
	switch {
	case errors.Is(err, BadRequest) || strings.Contains(err.Error(), ErrBadRequest):
		return NewRestError(http.StatusBadRequest, ErrBadRequest, err.Error(), debug)
	case errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), ErrNotFound):
		return NewRestError(http.StatusNotFound, ErrNotFound, err.Error(), debug)
	case errors.Is(err, context.DeadlineExceeded):
		return NewRestError(http.StatusRequestTimeout, ErrRequestTimeout, err.Error(), debug)
	case errors.Is(err, Unauthorized) || strings.Contains(err.Error(), ErrUnauthorized):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized, err.Error(), debug)
	case errors.Is(err, WrongCredentials) || strings.Contains(err.Error(), ErrWrongCredentials):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), "field validation"):
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return NewRestError(http.StatusBadRequest, ErrBadRequest, validationErrors.Error(), debug)
		}
		return parseValidatorError(err, debug)
	case strings.Contains(strings.ToLower(err.Error()), "required header"):
		return NewRestError(http.StatusBadRequest, ErrBadRequest, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), "no documents in result"):
		return NewRestError(http.StatusNotFound, ErrNotFound, err.Error(), debug)
	default:
		if restErr, ok := err.(*RestError); ok {
			return restErr
		}
		return NewRestError(http.StatusInternalServerError, ErrInternalServerError, errors.Cause(err).Error(), debug)
	}
}

func parseValidatorError(err error, debug bool) RestErr {
	if strings.Contains(err.Error(), "Password") {
		return NewRestError(http.StatusBadRequest, ErrInvalidPassword, err, debug)
	}

	if strings.Contains(err.Error(), "Email") {
		return NewRestError(http.StatusBadRequest, ErrInvalidEmail, err, debug)
	}

	return NewRestError(http.StatusBadRequest, ErrInvalidField, err, debug)
}

// ErrorResponse Error response
func ErrorResponse(err error, debug bool) (int, interface{}) {
	return ParseErrors(err, debug).Status(), ParseErrors(err, debug)
}

// ErrorCtxResponse Error response object and status code
func ErrorCtxResponse(ctx echo.Context, err error, debug bool) error {
	if err != nil {
		restErr := ParseErrors(err, debug)
		return ctx.JSON(restErr.Status(), restErr)
	}
	return ctx.JSON(http.StatusInternalServerError, ErrInternalServerError)
}
