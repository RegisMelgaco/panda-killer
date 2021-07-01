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

func (u TransferUsecase) CreateTransfer(ctx context.Context, newTransfer *transfer.Transfer) error {
	entry := logrus.WithField("transfer", newTransfer)

	originAccount, err := u.accountRepo.GetAccount(ctx, newTransfer.AccountOrigin)
	if err != nil {
		entry.Errorf("Failed to load originAccount on transfer creation with internal error: %v", err)
		return err
	}
	destinationAccount, err := u.accountRepo.GetAccount(ctx, newTransfer.AccountDestination)
	if err != nil {
		entry.Errorf("Failed to load destinationAccount on transfer creation with internal error: %v", err)
		return err
	}

	originNewBalance := safeSubtraction(originAccount.Balance, newTransfer.Amount)
	destinationNewBalance := safeSum(destinationAccount.Balance, newTransfer.Amount)

	err = u.transferRepo.CreateTransferAndUpdateAccountsBalances(ctx, newTransfer, originNewBalance, destinationNewBalance)
	if err != nil {
		logrus.New().WithField("transfer", newTransfer).Errorf("Failed to create transaction with internal error: %v", err)
		return err
	}

	return nil
}

func safeSubtraction(a, b float64) float64 {
	return float64(float64(int(a*100)-int(b*100)) / 100)
}

func safeSum(a, b float64) float64 {
	return float64(float64(int(a*100)+int(b*100)) / 100)
}
