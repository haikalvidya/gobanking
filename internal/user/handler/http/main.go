package http

import (
	"gobanking/internal/user/config"
	"gobanking/internal/user/usecase"
	"gobanking/pkg/logger"
	"gobanking/pkg/middlewares"
	"gobanking/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	User *UserHandler
	Auth *AuthHandler
}

type handler struct {
	Usecase   *usecase.Usecase
	Logger    logger.Logger
	cfg       *config.Config
	mw        middlewares.MiddlewareManager
	validator *utils.CustomValidator
	redis     *redis.Client
}

func NewHandler(usecase *usecase.Usecase,
	logger logger.Logger,
	cfg *config.Config,
	mw middlewares.MiddlewareManager,
	validator *utils.CustomValidator,
	e *echo.Group,
	redis *redis.Client) *Handler {
	handler := &handler{
		Usecase:   usecase,
		Logger:    logger,
		cfg:       cfg,
		mw:        mw,
		validator: validator,
		redis:     redis,
	}

	h := &Handler{
		User: (*UserHandler)(handler),
		Auth: (*AuthHandler)(handler),
	}

	auth := e.Group("/auth")
	{
		auth.POST("/signup", h.Auth.SignUp)
		auth.POST("/signin", h.Auth.SignIn)
		auth.POST("/refresh", h.Auth.RefreshToken)
		auth.POST("/signout", h.Auth.SignOut, handler.mw.AuthMiddleware)
		auth.GET("/me", h.Auth.Me, handler.mw.AuthMiddleware)
	}

	user := e.Group("/user")
	{
		user.GET("/:id", h.User.GetByID, handler.mw.AuthMiddleware)
		user.PUT("/:id", h.User.Update, handler.mw.AuthMiddleware)
		user.DELETE("/:id", h.User.Delete, handler.mw.AuthMiddleware)
	}

	return h
}
