package rest

import (
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/shared"
)

type CreateAccountRequest struct {
	Balance  shared.Money `json:"balance"`
	Name     string       `json:"name"`
	CPF      string       `json:"cpf" minLength:"11" maxLength:"11"`
	Password string       `json:"password"`
}

type CreateTransferRequest struct {
	OriginAccountID      account.AccountID `json:"origin_account_id"`
	DestinationAccountID account.AccountID `json:"destination_account_id"`
	Amount               shared.Money      `json:"amount"`
}

type LoginRequest struct {
	CPF      string `json:"cpf"`
	Password string `json:"password"`
}
