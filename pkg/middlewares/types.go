package middlewares

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password,omitempty"`
	FirsName  string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
