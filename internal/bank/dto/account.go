package dto

type CreateAccountRequest struct {
	Name       string `json:"name" validate:"required"`
	CurrencyId int    `json:"currency_id" validate:"required"`
}

type UpdateAccountRequest struct {
	Name string `json:"name"`
}
