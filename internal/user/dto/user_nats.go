package dto

import "gobanking/internal/user/models"

type GetUserByIDRequest struct {
	Token string `json:"id" validate:"required"`
}

type GetUserByIDResponse struct {
	Data models.User `json:"data"`
}
