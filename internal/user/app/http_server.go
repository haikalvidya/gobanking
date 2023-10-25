package app

import (
	httpError "gobanking/pkg/http_errors"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	maxHeaderBytes = 1 << 20
	stackSize      = 1 << 10 // 1 KB
	bodyLimit      = "2M"
	readTimeout    = 15 * time.Second
	writeTimeout   = 15 * time.Second
	gzipLevel      = 5
)

func (a *app) runHttpServer() error {
	a.settingMiddleware()

	a.echo.Server.ReadTimeout = readTimeout
	a.echo.Server.WriteTimeout = writeTimeout
	a.echo.Server.MaxHeaderBytes = maxHeaderBytes

	// add exceptional error handler
	a.echo.HTTPErrorHandler = a.errorHandler

	return a.echo.Start(a.cfg.Http.Port)
}

func (a *app) errorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := http.StatusText(code)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message.(string)
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else if code == http.StatusNotFound {
			// if not found, return 404
			err = httpError.NewNotFoundError(c, msg, a.cfg.Http.DebugErrorsResponse)
		} else if code == http.StatusMethodNotAllowed {
			// if method not allowed, return 405
			err = httpError.NewMethodNotAllowedError(c, msg, a.cfg.Http.DebugErrorsResponse)
		} else {
			err = httpError.NewInternalServerError(c, msg, a.cfg.Http.DebugErrorsResponse)
		}
		if err != nil {
			a.log.Errorf("(errorHandler) err: %v", err)
		}
	}
}

func (a *app) settingMiddleware() {
	a.echo.Use(a.middlewareManager.RequestLoggerMiddleware)
	a.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowCredentials: true,
	}))
	a.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         stackSize,
		DisablePrintStack: false,
		DisableStackAll:   false,
	}))
	a.echo.Use(middleware.RequestID())
	a.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: gzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	a.echo.Use(middleware.BodyLimit(bodyLimit))
}
