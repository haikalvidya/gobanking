package middlewares

import (
	dtoUser "gobanking/internal/user/dto"
	"gobanking/pkg/http_errors"
	natsPkg "gobanking/pkg/nats"
	"gobanking/pkg/serializer"
	"time"

	"github.com/labstack/echo/v4"
)

type UserResponse struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	FirsName  string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// this middleware is for the service that will call user service
func (mw *middlewareManager) AuthMiddlewareClient(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// get jwt token from header
		token := ctx.Request().Header.Get("Authorization")

		req := dtoUser.GetUserByIDRequest{}
		req.ID = token

		reqBytes, err := serializer.Marshal(req)
		if err != nil {
			mw.log.Errorf("failed to marshal request: %v", err)
			return http_errors.NewUnauthorizedError(ctx, err, mw.config.DebugErrorsResponse)
		}

		// call user api to validate token to nats
		resp, err := mw.natsConn.Request(natsPkg.UserGetUserByIdReqRep, reqBytes, natsPkg.TimeoutReq)
		if err != nil {
			mw.log.Errorf("failed to request to nats: %v", err)
			return http_errors.NewUnauthorizedError(ctx, err, mw.config.DebugErrorsResponse)
		}

		// get user data from response
		var userResponse UserResponse
		if err := serializer.Unmarshal(resp.Data, &userResponse); err != nil {
			mw.log.Errorf("failed to unmarshal response: %v", err)
			return http_errors.NewUnauthorizedError(ctx, err, mw.config.DebugErrorsResponse)
		}

		// set user data to context
		ctx.Set("user", userResponse)

		return next(ctx)
	}
}
