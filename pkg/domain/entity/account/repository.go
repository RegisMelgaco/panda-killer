package account

import "context"

type AccountRepo interface {
	CreateAccount(context.Context, *Account) error
	GetAccounts(context.Context) ([]*Account, error)
	GetAccountBalance(context.Context, int) (float64, error)
}
