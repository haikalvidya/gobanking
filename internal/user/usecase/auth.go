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

	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	// SingUp creates a new user
	SignUp(ctx context.Context, user *dto.SignUpRequest) (*dto.SignUpResponse, error)
	// SignIn signs in a user
	SignIn(ctx context.Context, user *dto.SignInRequest) (*dto.SignInResponse, error)
	// RefreshToken refreshes a token
	RefreshToken(ctx context.Context, token string) (*dto.RefreshTokenResponse, error)
	// SignOut signs out a user
	SignOut(ctx context.Context, token string) error
	// Me returns user data
	Me(ctx context.Context, token string) (*models.User, error)
}

type authUsecase usecase

type claim struct {
	UserId    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	jwt.StandardClaims
}

func (u *authUsecase) generateToken(ctx context.Context, user *models.User) (string, string, error) {
	exp := time.Until(time.Unix(u.cfg.JWT.Expire, 0))
	claimsToken := claim{
		UserId:    user.ID.String(),
		UserEmail: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Milliseconds(),
			Issuer:    u.cfg.JWT.Issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	theSecret := []byte(u.cfg.JWT.Secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsToken)
	tokenString, err := token.SignedString(theSecret)
	if err != nil {
		return "", "", fmt.Errorf("error signing token: %v", err)
	}

	expRefresh := time.Until(time.Unix(u.cfg.JWT.RefreshExpire, 0))

	claimsRefreshToken := claim{
		UserId:    user.ID.String(),
		UserEmail: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expRefresh.Milliseconds(),
			Issuer:    u.cfg.JWT.Issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefreshToken)
	refreshTokenString, err := refreshToken.SignedString(theSecret)
	if err != nil {
		return "", "", fmt.Errorf("error signing refresh token: %v", err)
	}

	return tokenString, refreshTokenString, nil
}

func (u *authUsecase) saveSessionTokenToRedis(ctx context.Context, token, refreshToken string, user *models.User) error {
	// serialize user
	userJson, err := serializer.Marshal(user)
	if err != nil {
		return fmt.Errorf("error serializing user: %v", err)
	}

	keyToken := redisPkg.GetKeyOfTokenUserFromRedis(token, redisPkg.Token)
	exp := time.Until(time.Unix(u.cfg.JWT.Expire, 0))
	if err := u.Redis.Set(ctx, keyToken, userJson, exp).Err(); err != nil {
		return fmt.Errorf("error saving token to redis: %v", err)
	}

	keyRefreshToken := redisPkg.GetKeyOfTokenUserFromRedis(refreshToken, redisPkg.RefreshToken)
	expRefresh := time.Until(time.Unix(u.cfg.JWT.RefreshExpire, 0))
	if err := u.Redis.Set(ctx, keyRefreshToken, userJson, expRefresh).Err(); err != nil {
		return fmt.Errorf("error saving refresh token to redis: %v", err)
	}

	return nil
}

func (u *authUsecase) SignUp(ctx context.Context, user *dto.SignUpRequest) (*dto.SignUpResponse, error) {
	// check if user already exists by email and username
	userFound, err := u.Repo.User.IsUserExist(ctx, user.Email, user.Username)
	if err != nil {
		u.Logger.Errorf("error checking if user exists: %v", err)
		typeOfError := http_errors.BadRequest
		return nil, fmt.Errorf("%v : error checking if user exists: %v", typeOfError, err)
	}

	if userFound {
		u.Logger.Errorf("user already exists")
		typeOfError := http_errors.BadRequest
		return nil, fmt.Errorf("%v : user already exists", typeOfError)
	}

	// create user
	userModel := &models.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		FirsName: user.FirstName,
		LastName: user.LastName,
	}

	if err := userModel.HashPassword(); err != nil {
		u.Logger.Errorf("error hashing password: %v", err)
		typeOfError := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error hashing password: %v", typeOfError, err)
	}

	if err := u.Repo.User.Create(ctx, userModel); err != nil {
		u.Logger.Errorf("error creating user: %v", err)
		typeOfError := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error creating user: %v", typeOfError, err)
	}

	// get user from db
	userModel, err = u.Repo.User.GetByEmail(ctx, user.Email)
	if err != nil {
		u.Logger.Errorf("error getting user by email: %v", err)
		typeOfError := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error getting user by email: %v", typeOfError, err)
	}

	token, refreshToken, err := u.generateToken(ctx, userModel)
	if err != nil {
		u.Logger.Errorf("error generating token: %v", err)
		typeOfError := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error generating token: %v", typeOfError, err)
	}

	// save token to redis
	if err := u.saveSessionTokenToRedis(ctx, token, refreshToken, userModel); err != nil {
		u.Logger.Errorf("error saving token to redis: %v", err)
		typeOfError := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error saving token to redis: %v", typeOfError, err)
	}

	return &dto.SignUpResponse{
		ID:           userModel.ID.String(),
		Username:     userModel.Username,
		Email:        userModel.Email,
		FirsName:     userModel.FirsName,
		LastName:     userModel.LastName,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (u *authUsecase) SignIn(ctx context.Context, user *dto.SignInRequest) (*dto.SignInResponse, error) {
	// get by username or email
	userModel, err := u.Repo.User.GetByEmail(ctx, user.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Logger.Errorf("error getting user by email: %v", err)
		typeOfError := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error getting user by email: %v", typeOfError, err)
	}
	if err == gorm.ErrRecordNotFound {
		u.Logger.Errorf("user not found")
		typeOfError := http_errors.BadRequest
		return nil, fmt.Errorf("%v : user not found", typeOfError)
	}

	// check password
	if err := userModel.ComparePasswords(user.Password); err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
		u.Logger.Errorf("error checking password: %v", err)
		typeOfError := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error checking password: %v", typeOfError, err)
	}

	if err == bcrypt.ErrMismatchedHashAndPassword {
		u.Logger.Errorf("invalid password")
		typeOfError := http_errors.BadRequest
		return nil, fmt.Errorf("%v : invalid password", typeOfError)
	}

	// generate token
	token, refreshToken, err := u.generateToken(ctx, userModel)
	if err != nil {
		u.Logger.Errorf("error generating token: %v", err)
		typeOfError := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error generating token: %v", typeOfError, err)
	}

	// save token to redis
	if err := u.saveSessionTokenToRedis(ctx, token, refreshToken, userModel); err != nil {
		u.Logger.Errorf("error saving token to redis: %v", err)
		typeOfError := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error saving token to redis: %v", typeOfError, err)
	}

	return &dto.SignInResponse{
		ID:           userModel.ID.String(),
		Username:     userModel.Username,
		Email:        userModel.Email,
		FirsName:     userModel.FirsName,
		LastName:     userModel.LastName,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (u *authUsecase) RefreshToken(ctx context.Context, token string) (*dto.RefreshTokenResponse, error) {
	var err error
	// check if token is in redis
	key := redisPkg.GetKeyOfTokenUserFromRedis(token, redisPkg.RefreshToken)
	var user *models.User
	user, err = redisPkg.GetDataFromRedis[models.User](ctx, u.Redis, key, user)
	if err != nil && err != redis.Nil {
		u.Logger.Errorf("error getting redis user: %v", err)
		typeOfError := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error getting redis user: %v", typeOfError, err)
	}

	if err == redis.Nil {
		u.Logger.Errorf("Unauthorized")
		typeOfError := http_errors.Unauthorized
		return nil, fmt.Errorf(typeOfError.Error())
	}

	// delete refresh from redis
	if err := u.Redis.Del(ctx, key).Err(); err != nil {
		u.Logger.Errorf("error deleting refresh token from redis: %v", err)
		typeOfError := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error deleting refresh token from redis: %v", typeOfError, err)
	}

	// generate token
	newToken, newRefreshToken, err := u.generateToken(ctx, user)
	if err != nil {
		u.Logger.Errorf("error generating token: %v", err)
		typeOfError := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error generating token: %v", typeOfError, err)
	}

	// save token to redis
	if err := u.saveSessionTokenToRedis(ctx, newToken, newRefreshToken, user); err != nil {
		u.Logger.Errorf("error saving token to redis: %v", err)
		typeOfError := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error saving token to redis: %v", typeOfError, err)
	}

	return &dto.RefreshTokenResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (u *authUsecase) SignOut(ctx context.Context, token string) error {
	// delete token from redis
	key := redisPkg.GetKeyOfTokenUserFromRedis(token, redisPkg.Token)

	if err := u.Redis.Del(ctx, key).Err(); err != nil {
		u.Logger.Errorf("error deleting token from redis: %v", err)
		typeOfError := http_errors.InternalServerError
		return fmt.Errorf("%v : error deleting token from redis: %v", typeOfError, err)
	}

	return nil
}

func (u *authUsecase) Me(ctx context.Context, token string) (*models.User, error) {
	// get user from redis
	key := redisPkg.GetKeyOfTokenUserFromRedis(token, redisPkg.Token)

	var user *models.User

	user, err := redisPkg.GetDataFromRedis[models.User](ctx, u.Redis, key, user)
	if err != nil && err != redis.Nil {
		u.Logger.Errorf("error getting redis user: %v", err)
		typeOfError := http_errors.InternalServerError
		return nil, fmt.Errorf("%v : error getting redis user: %v", typeOfError, err)
	}

	if err == redis.Nil {
		u.Logger.Errorf("token is invalid")
		typeOfError := http_errors.BadRequest
		return nil, fmt.Errorf("%v : token is invalid", typeOfError)
	}

	user.Clean()

	return user, nil
}
