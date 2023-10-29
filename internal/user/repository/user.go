package repository

import (
	"context"
	"gobanking/internal/user/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	// Create creates a new user
	CreateTX(ctx context.Context, tx *gorm.DB, user *models.User) error
	// GetByEmail returns a user by email
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	// GetByUsername returns a user by username
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	// IsUserExist checks if a user exists by email or username
	IsUserExist(ctx context.Context, email, username string) (bool, error)
	// GetByID returns a user by id
	GetByID(ctx context.Context, id string) (*models.User, error)
	// Update updates a user
	UpdateTX(ctx context.Context, tx *gorm.DB, user *models.User) error
	// Delete deletes a user
	Delete(ctx context.Context, id string) error
}

type userRepository repository

func (u *userRepository) CreateTX(ctx context.Context, tx *gorm.DB, user *models.User) error {
	return tx.WithContext(ctx).Create(user).Error
}

func (u *userRepository) IsUserExist(ctx context.Context, email, username string) (bool, error) {
	var count int64
	err := u.DB.WithContext(ctx).Model(&models.User{}).Where("email = ? OR username = ?", email, username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := u.DB.WithContext(ctx).Where("email = ?", email).First(user).Error
	return user, err
}

func (u *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	err := u.DB.WithContext(ctx).Where("username = ?", username).First(user).Error
	return user, err
}

func (u *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	err := u.DB.WithContext(ctx).Where("id = ?", id).First(user).Error
	return user, err
}

func (u *userRepository) UpdateTX(ctx context.Context, tx *gorm.DB, user *models.User) error {
	return tx.WithContext(ctx).Save(user).Error
}

func (u *userRepository) Delete(ctx context.Context, id string) error {
	// get user by id
	user, err := u.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// delete user
	return u.DB.WithContext(ctx).Delete(user).Error
}
