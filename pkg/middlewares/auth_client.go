package middlewares

import (
	"github.com/labstack/echo/v4"
)

// this middleware is for the service that will call user service
func (mw *middlewareManager) AuthMiddlewareClient(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// // get jwt token from header
		// token := ctx.Request().Header.Get("Authorization")

		// // call user api to validate token
		// resp, err := mw.httpClient.R().
		// 	SetHeader("Authorization", token).
		// 	Get(mw.config.UserApiUrl + "/mine")
		// if err != nil {
		// 	mw.log.Error(err)
		// 	return httpErrors.NewInternalServerError(ctx, err, mw.config.DebugErrorsResponse)
		// }

		// if resp.StatusCode() != http.StatusOK {
		// 	return httpErrors.NewUnauthorizedError(ctx, "Unauthorized", mw.config.DebugErrorsResponse)
		// }

		// // get user data from response
		// var userApiResponse UserApiResponse
		// if err := serializer.Unmarshal(resp.Body(), &userApiResponse); err != nil {
		// 	mw.log.Error(err)
		// 	return httpErrors.NewUnauthorizedError(ctx, err, mw.config.DebugErrorsResponse)
		// }

		// // set user data to context
		// ctx.Set("user", userApiResponse)

		return next(ctx)
	}
}
