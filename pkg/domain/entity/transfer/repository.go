package transfer

import (
	"context"
	"local/panda-killer/pkg/domain/entity/account"
)

type TransferRepo interface {
	CreateTransferAndUpdateAccountsBalances(context.Context, *Transfer) error
	GetTransfersCantainingAccount(context.Context, account.AccountID) ([]Transfer, error)
}
