package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
)

type UserApiResponse struct {
	Data struct {
		ID        int         `json:"id"`
		Email     string      `json:"email"`
		Username  string      `json:"username"`
		Name      string      `json:"name"`
		Phone     string      `json:"phone"`
		PhotoURL  interface{} `json:"photoUrl"`
		BranchID  int         `json:"branch_id"`
		RoleID    int         `json:"role_id"`
		IsDeleted int         `json:"is_deleted"`
		CreatedAt time.Time   `json:"created_at"`
		UpdatedAt interface{} `json:"updated_at"`
		DeletedAt interface{} `json:"deleted_at"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (mw *middlewareManager) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// get jwt token from header
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
