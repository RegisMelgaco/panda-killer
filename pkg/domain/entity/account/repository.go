package account

import "context"

type AccountRepo interface {
	CreateAccount(context.Context, *Account) error
	GetAccounts(context.Context) ([]*Account, error)
	GetAccount(context.Context, int) (*Account, error)
	GetAccountByCPF(context.Context, string) (*Account, error)
}
