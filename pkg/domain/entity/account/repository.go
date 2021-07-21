package account

import "context"

type AccountRepo interface {
	CreateAccount(context.Context, *Account) error
	GetAccounts(context.Context) ([]Account, error)
	GetAccount(context.Context, AccountID) (*Account, error)
	GetAccountByCPF(ctx context.Context, cpf string) (*Account, error)
}
