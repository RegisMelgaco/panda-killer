package rest

import "time"

type AccountBalanceResponse struct {
	Balance int `json:"balance"`
}

type ErrorResponse struct {
	Message string
}

type CreatedAccountResponse struct {
	ID int `json:"id"`
}

type GetAccountResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	CPF  string `json:"cpf"`
}

type CreateTransferResponse struct {
	ID int `json:"id"`
}

type GetTransferResponse struct {
	ID                   int       `json:"id"`
	Amount               int       `json:"amount"`
	OriginAccountID      int       `json:"origin_account_id"`
	DestinationAccountID int       `json:"destination_account_id"`
	CreatedAt            time.Time `json:"created_at"`
}
