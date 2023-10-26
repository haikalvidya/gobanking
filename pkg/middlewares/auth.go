package middlewares

import (
	"gobanking/pkg/http_errors"
	redisPkg "gobanking/pkg/redis"
	"html"

	"github.com/labstack/echo/v4"
)

func (mw *middlewareManager) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// get jwt token from header
		bearerToken := ctx.Request().Header.Get("Authorization")
		// check if token is valid
		if bearerToken == "" {
			mw.log.Error("token is empty")
			return http_errors.NewUnauthorizedError(ctx, "Unauthorized", mw.config.DebugErrorsResponse)
		}
		token := html.EscapeString(bearerToken)

		// to struct
		var user *UserModel
		key := redisPkg.GetKeyOfTokenUserFromRedis(token, redisPkg.Token)

		user, err := redisPkg.GetDataFromRedis(ctx.Request().Context(), mw.redis, key, user)
		if err != nil {
			mw.log.Errorf("failed to get redis user otp: %v", err)
			return http_errors.NewUnauthorizedError(ctx, "Unauthorized", mw.config.DebugErrorsResponse)
		}
		// set user data to context
		ctx.Set("user", user)

		return next(ctx)
	}
}
