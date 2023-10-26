package middlewares

import (
	"gobanking/pkg/logger"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type MiddlewareManager interface {
	RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type middlewareManager struct {
	log    logger.Logger
	config *MiddlewareConfig
	redis  *redis.Client
}

type MiddlewareConfig struct {
	HttpClientDebug     bool
	DebugErrorsResponse bool
}

func NewMiddlewareManager(log logger.Logger, cfg *MiddlewareConfig, redis *redis.Client) *middlewareManager {
	mwManager := &middlewareManager{log: log, config: cfg, redis: redis}
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
