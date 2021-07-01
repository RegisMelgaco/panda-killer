package usecase

import (
	"context"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/transfer"

	"github.com/sirupsen/logrus"
)

type TransferUsecase struct {
	transferRepo transfer.TransferRepo
	accountRepo  account.AccountRepo
}

func NewTransferUsecase(transferRepo transfer.TransferRepo, accountRepo account.AccountRepo) *TransferUsecase {
	return &TransferUsecase{transferRepo: transferRepo, accountRepo: accountRepo}
}

func (u TransferUsecase) CreateTransfer(ctx context.Context, originAccountID, destinationAccountID int, amount float64) (transfer.Transfer, error) {
	entry := logrus.WithFields(logrus.Fields{
		"originAccountID":      originAccountID,
		"destinationAccountID": destinationAccountID,
		"amount":               amount,
	})

	originAccount, err := u.accountRepo.GetAccount(ctx, originAccountID)
	if err != nil {
		entry.Errorf("Failed to load originAccount on transfer creation with internal error: %v", err)
		return transfer.Transfer{}, err
	}
	destinationAccount, err := u.accountRepo.GetAccount(ctx, destinationAccountID)
	if err != nil {
		entry.Errorf("Failed to load destinationAccount on transfer creation with internal error: %v", err)
		return transfer.Transfer{}, err
	}

	newTransfer := transfer.NewTransfer(originAccount, destinationAccount, amount)

	err = u.transferRepo.CreateTransferAndUpdateAccountsBalances(ctx, newTransfer)
	if err != nil {
		logrus.New().WithField("transfer", newTransfer).Errorf("Failed to create transaction with internal error: %v", err)
		return transfer.Transfer{}, err
	}

	return *newTransfer, nil
}
