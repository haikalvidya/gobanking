package dto

type TransferRequest struct {
	SourceAccountNumber      string `json:"source_account_number" validate:"required"`
	DestinationAccountNumber string `json:"destination_account_number" validate:"required"`
	Amount                   int64  `json:"amount" validate:"required"`
}

type WithdrawalRequest struct {
	SourceAccountNumber string `json:"source_account_number" validate:"required"`
	Amount              int64  `json:"amount" validate:"required"`
}

type DepositRequest struct {
	DestinationAccountNumber string `json:"destination_account_number" validate:"required"`
	Amount                   int64  `json:"amount" validate:"required"`
}
