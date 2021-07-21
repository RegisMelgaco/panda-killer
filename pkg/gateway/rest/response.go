package rest

import (
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/shared"
	"local/panda-killer/pkg/domain/entity/transfer"
	"time"
)

type AccountBalanceResponse struct {
	Balance shared.Money `json:"balance"`
}

type ErrorResponse struct {
	Message string
}

type CreatedAccountResponse struct {
	ID account.AccountID `json:"id"`
}

type GetAccountResponse struct {
	ID   account.AccountID `json:"id"`
	Name string            `json:"name"`
	CPF  string            `json:"cpf"`
}

type CreateTransferResponse struct {
	ID transfer.TransferID `json:"id"`
}

type GetTransferResponse struct {
	ID                   transfer.TransferID `json:"id"`
	Amount               shared.Money        `json:"amount"`
	OriginAccountID      account.AccountID   `json:"origin_account_id"`
	DestinationAccountID account.AccountID   `json:"destination_account_id"`
	CreatedAt            time.Time           `json:"created_at"`
}
