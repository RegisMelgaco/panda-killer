package transfer

import (
	"local/panda-killer/pkg/domain/entity/account"
	"time"
)

type Transfer struct {
	ID                                int
	OriginAccount, DestinationAccount *account.Account
	Amount                            float64
	CreatedAt                         time.Time
}

func NewTransfer(originAccount, destinationAccount *account.Account, amount float64) (*Transfer, error) {
	if amount <= 0 {
		return &Transfer{}, ErrTransferAmountShouldBeGreatterThanZero
	}
	if originAccount.Balance < amount {
		return &Transfer{}, ErrInsufficientFundsToMakeTransaction
	}

	originAccount.Balance = safeSubtraction(originAccount.Balance, amount)
	destinationAccount.Balance = safeSum(destinationAccount.Balance, amount)
	return &Transfer{
		OriginAccount:      originAccount,
		DestinationAccount: destinationAccount,
		Amount:             amount,
		CreatedAt:          time.Now(),
	}, nil
}

func safeSubtraction(a, b float64) float64 {
	return float64(float64(int(a*100)-int(b*100)) / 100)
}

func safeSum(a, b float64) float64 {
	return float64(float64(int(a*100)+int(b*100)) / 100)
}
