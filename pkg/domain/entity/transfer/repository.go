package transfer

import "context"

type TransferRepo interface {
	CreateTransferAndUpdateAccountsBalances(context.Context, *Transfer, float64, float64) error
}
