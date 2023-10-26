package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	User UserRepository
	Tx   Tx
}

type Tx interface {
	DoInTransaction(fn func(tx *gorm.DB) error) (err error)
}

type tx struct {
	DB *gorm.DB
}

func (t *tx) DoInTransaction(fn func(tx *gorm.DB) error) (err error) {

	tx := t.DB.Begin()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	err = fn(tx)

	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	return
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
		Tx:   &tx{DB: db},
	}
}
