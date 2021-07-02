package transfer

import (
	"local/panda-killer/pkg/domain/entity/account"
	"time"
)

type Transfer struct {
	ID                                int
	OriginAccount, DestinationAccount *account.Account
	Amount                            int
	CreatedAt                         time.Time
}

func NewTransfer(originAccount, destinationAccount *account.Account, amount int) (*Transfer, error) {
	if originAccount.ID == destinationAccount.ID {
		return &Transfer{}, ErrTransferOriginAndDestinationNeedToBeDiffrent
	}
	if amount <= 0 {
		return &Transfer{}, ErrTransferAmountShouldBeGreatterThanZero
	}
	if originAccount.Balance < amount {
		return &Transfer{}, ErrInsufficientFundsToMakeTransaction
	}

	originAccount.Balance = originAccount.Balance - amount
	destinationAccount.Balance = destinationAccount.Balance + amount
	return &Transfer{
		OriginAccount:      originAccount,
		DestinationAccount: destinationAccount,
		Amount:             amount,
		CreatedAt:          time.Now(),
	}, nil
}
