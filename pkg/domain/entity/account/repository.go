package account

import "context"

//go:generate moq -out repository_mock.go . AccountRepo

type AccountRepo interface {
	CreateAccount(context.Context, *Account) error
	GetAccounts(context.Context) ([]Account, error)
	GetAccount(context.Context, AccountID) (*Account, error)
	GetAccountByCPF(ctx context.Context, cpf string) (*Account, error)
}
