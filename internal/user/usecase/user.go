package usecase

import (
	"context"
	"fmt"
	"gobanking/internal/user/dto"
	"gobanking/internal/user/models"
	"gobanking/pkg/http_errors"
	redisPkg "gobanking/pkg/redis"
	"gobanking/pkg/serializer"
	"time"

	"gorm.io/gorm"
)

type UserUsecase interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
	// Update updates a user
	Update(ctx context.Context, id, token string, user *dto.UpdateUserRequest) (*models.User, error)
	// Delete deletes a user
	Delete(ctx context.Context, id, token string) error
}

type userUsecase usecase

func (u *userUsecase) updateUserRedis(ctx context.Context, token string, user *models.User) error {
	// delete user from redis first
	keyToken := redisPkg.GetKeyOfTokenUserFromRedis(token, redisPkg.Token)
	if err := u.Redis.Del(ctx, keyToken).Err(); err != nil {
		return fmt.Errorf("error deleting refresh token from redis")
	}

	// serialize user
	userJson, err := serializer.Marshal(user)
	if err != nil {
		return fmt.Errorf("error serializing user")
	}

	exp := time.Until(time.Unix(u.cfg.JWT.Expire, 0))
	// set user to redis
	if err := u.Redis.Set(ctx, keyToken, userJson, exp).Err(); err != nil {
		return fmt.Errorf("error setting user to redis")
	}

	return nil
}

func (u *userUsecase) GetByID(ctx context.Context, id string) (*models.User, error) {
	// get user from db
	user, err := u.Repo.User.GetByID(ctx, id)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Logger.Errorf("error when get user by id: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when get user by id: %v", typeOfErr, err)
	}

	if err == gorm.ErrRecordNotFound {
		u.Logger.Errorf("user not found")
		typeOfErr := http_errors.NotFound
		return nil, fmt.Errorf("%v : user not found", typeOfErr)
	}

	user.Clean()

	return user, nil
}

func (u *userUsecase) Update(ctx context.Context, id, token string, user *dto.UpdateUserRequest) (*models.User, error) {
	// get user from db by id
	userModel, err := u.Repo.User.GetByID(ctx, id)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Logger.Errorf("error when get user by id: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when get user by id: %v", typeOfErr, err)
	}

	if err == gorm.ErrRecordNotFound {
		u.Logger.Errorf("user not found")
		typeOfErr := http_errors.NotFound
		return nil, fmt.Errorf("%v : user not found", typeOfErr)
	}

	// update user if there is a change
	if user.FirstName != "" {
		userModel.FirsName = user.FirstName
	}
	if user.LastName != "" {
		userModel.LastName = user.LastName
	}
	if user.Password != "" {
		userModel.Password = user.Password
		// hash password
		if err := userModel.HashPassword(); err != nil {
			u.Logger.Errorf("error when hash password: %v", err)
			typeOfErr := http_errors.InternalServerError
			return nil, fmt.Errorf("%v : error when hash password: %v", typeOfErr, err)
		}
	}

	// update user in db
	if err := u.Repo.User.Update(ctx, userModel); err != nil {
		u.Logger.Errorf("error when update user: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when update user: %v", typeOfErr, err)
	}

	userModel.Clean()

	// update user in redis
	if err := u.updateUserRedis(ctx, token, userModel); err != nil {
		u.Logger.Errorf("error when update user in redis: %v", err)
		typeOfErr := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error when update user in redis: %v", typeOfErr, err)
	}

	return userModel, nil
}

func (u *userUsecase) Delete(ctx context.Context, id, token string) error {
	// delete user from db
	err := u.Repo.User.Delete(ctx, id)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Logger.Errorf("error when delete user: %v", err)
		typeOfErr := http_errors.InternalServerError
		return fmt.Errorf("%v : error when delete user: %v", typeOfErr, err)
	}

	if err == gorm.ErrRecordNotFound {
		u.Logger.Errorf("user not found")
		typeOfErr := http_errors.NotFound
		return fmt.Errorf("%v : user not found", typeOfErr)
	}

	// delete user from redis
	keyToken := redisPkg.GetKeyOfTokenUserFromRedis(token, redisPkg.Token)
	if err := u.Redis.Del(ctx, keyToken).Err(); err != nil {
		u.Logger.Errorf("error deleting refresh token from redis")
		typeOfErr := http_errors.InternalServerError
		return fmt.Errorf("%v : error deleting refresh token from redis", typeOfErr)
	}

	return nil
}
