package middlewares

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	ID        uuid.UUID  `json:"id" gorm:"primary_key"`
	Username  string     `json:"username" gorm:"unique;not null"`
	Email     string     `json:"email" gorm:"unique;not null"`
	Password  string     `json:"password,omitempty" gorm:"not null"`
	FirsName  string     `json:"first_name" gorm:"not null"`
	LastName  string     `json:"last_name"`
	CreatedAt *time.Time `gorm:":autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `gorm:":autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:null"`
}
