package middlewares

import (
	"gobanking/pkg/logger"
	"time"

	"github.com/labstack/echo/v4"
)

type MiddlewareManager interface {
	RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type middlewareManager struct {
	log    logger.Logger
	config *MiddlewareConfig
}

type MiddlewareConfig struct {
	HttpClientDebug     bool
	DebugErrorsResponse bool
	UserApiUrl          string
}

func NewMiddlewareManager(log logger.Logger, cfg *MiddlewareConfig) *middlewareManager {
	mwManager := &middlewareManager{log: log, config: cfg}
	return mwManager
}

func (mw *middlewareManager) RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		start := time.Now()
		err := next(ctx)

		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		size := res.Size
		s := time.Since(start)

		mw.log.HttpMiddlewareAccessLogger(req.Method, req.URL.String(), status, size, s)

		return err
	}
}
