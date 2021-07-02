package transfer

import "context"

type TransferRepo interface {
	CreateTransferAndUpdateAccountsBalances(context.Context, *Transfer) error
	GetTransfersCantainingAccount(context.Context, int) ([]*Transfer, error)
}
