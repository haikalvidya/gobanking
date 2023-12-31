package utils

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

type FieldError struct {
	Param   string
	Message string
}

func NewFieldError(validator *validator.Validate) *CustomValidator {
	return &CustomValidator{Validator: validator}
}

func msgForTag(fe validator.FieldError, fieldName string) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s field is required", fieldName)
	case "numeric":
		return fmt.Sprintf("%s must be a number", fieldName)
	case "email":
		return "Email invalid"
	}
	return fe.Error()
}

func (cv *CustomValidator) Validate(i interface{}) error {
	cv.Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		}

		return name
	})
	if err := cv.Validator.Struct(i); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorRes := make(map[string]interface{})
			for _, fe := range ve {
				errorRes[fe.Field()] = msgForTag(fe, fe.Field())
			}
			return echo.NewHTTPError(http.StatusBadRequest, errorRes)
		}
	}
	return nil
}

func GetErrorValidation(err error) interface{} {
	return err.(*echo.HTTPError).Message
}
