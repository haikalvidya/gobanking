package usecase

import (
	"gobanking/internal/user/config"
	"gobanking/internal/user/repository"
	"gobanking/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type Usecase struct {
	User UserUsecase
	Auth AuthUsecase
}

type usecase struct {
	Repo   *repository.Repository
	Redis  *redis.Client
	Logger logger.Logger
	cfg    *config.Config
}

func NewUsecase(repo *repository.Repository, redis *redis.Client, logger logger.Logger, cfg *config.Config) *Usecase {
	svc := &usecase{
		Repo:   repo,
		Redis:  redis,
		Logger: logger,
		cfg:    cfg,
	}
	return &Usecase{
		User: (*userUsecase)(svc),
		Auth: (*authUsecase)(svc),
	}
}
