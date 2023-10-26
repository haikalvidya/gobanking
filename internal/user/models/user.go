package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"primary_key"`
	Username  string         `json:"username" gorm:"not null"`
	Email     string         `json:"email" gorm:"not null"`
	Password  string         `json:"password,omitempty" gorm:"not null"`
	FirsName  string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name"`
	CreatedAt *time.Time     `gorm:":autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time     `gorm:":autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"default:null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}

func (u *User) Clean() {
	u.Password = ""
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePasswords(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
