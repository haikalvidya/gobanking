package dto

import "gobanking/internal/user/models"

type GetUserByIDRequest struct {
	ID string `json:"id" validate:"required"`
}

type GetUserByIDResponse struct {
	Data models.User `json:"data"`
}
