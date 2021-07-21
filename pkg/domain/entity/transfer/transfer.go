package transfer

import (
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/shared"
	"time"
)

type Transfer struct {
	ID                                TransferID
	OriginAccount, DestinationAccount *account.Account
	Amount                            shared.Money
	CreatedAt                         time.Time
}

func NewTransfer(originAccount, destinationAccount *account.Account, amount shared.Money) (*Transfer, error) {
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
		Amount:             shared.Money(amount),
		CreatedAt:          time.Now(),
	}, nil
}
