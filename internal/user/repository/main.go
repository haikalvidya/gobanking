package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	Tx   Tx
	User UserRepository
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
		Tx:   &tx{DB: db},
		User: (*userRepository)(repo),
	}
}
