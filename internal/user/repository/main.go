package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	User UserRepository
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	repo := &repository{
		DB: db,
	}

	return &Repository{
		User: (*userRepository)(repo),
	}
}
