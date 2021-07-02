package usecase

import (
	"context"
	"fmt"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/transfer"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
)

type TransferUsecase struct {
	transferRepo transfer.TransferRepo
	accountRepo  account.AccountRepo
	group        *singleflight.Group
	mux          *sync.Mutex
}

func NewTransferUsecase(transferRepo transfer.TransferRepo, accountRepo account.AccountRepo) *TransferUsecase {
	return &TransferUsecase{
		transferRepo: transferRepo,
		accountRepo:  accountRepo,

		group: &singleflight.Group{},
		mux:   &sync.Mutex{},
	}
}

type createTransferFunc func(ctx context.Context, originAccountID, destinationAccountID, amount int) (*transfer.Transfer, error)

func (u TransferUsecase) CreateTransfer(ctx context.Context, originAccountID, destinationAccountID, amount int) (*transfer.Transfer, error) {
	return u.handleCreateTransferParallelism(ctx, originAccountID, destinationAccountID, amount, u.handleCreateTransferBussLogic)
}

func (u TransferUsecase) handleCreateTransferParallelism(ctx context.Context, originAccountID, destinationAccountID, amount int, f createTransferFunc) (*transfer.Transfer, error) {
	u.mux.Lock()

	// Prevents a dead lock with it self where is called Do two times on same key by same goroutine.
	if originAccountID == destinationAccountID {
		return &transfer.Transfer{}, transfer.ErrTransferOriginAndDestinationNeedToBeDiffrent
	}

	t, err, _ := u.group.Do(fmt.Sprint(originAccountID), func() (interface{}, error) {
		defer u.group.Forget(fmt.Sprint(originAccountID))
		t, err, _ := u.group.Do(fmt.Sprint(destinationAccountID), func() (interface{}, error) {
			defer u.group.Forget(fmt.Sprint(destinationAccountID))
			defer u.mux.Unlock()

			return f(ctx, originAccountID, destinationAccountID, amount)
		})
		return t, err
	})

	return t.(*transfer.Transfer), err
}

func (u TransferUsecase) handleCreateTransferBussLogic(ctx context.Context, originAccountID, destinationAccountID, amount int) (*transfer.Transfer, error) {
	entry := logrus.WithFields(logrus.Fields{
		"originAccountID":      originAccountID,
		"destinationAccountID": destinationAccountID,
		"amount":               amount,
	})

	originAccount, err := u.accountRepo.GetAccount(ctx, originAccountID)
	if err != nil {
		entry.Errorf("Failed to load originAccount on transfer creation with internal error: %v", err)
		return &transfer.Transfer{}, err
	}
	if originAccount.ID == 0 {
		entry.Warnf("Failed to create account with nonexisting originAccount")
		return &transfer.Transfer{}, account.ErrAccountNotFound
	}
	destinationAccount, err := u.accountRepo.GetAccount(ctx, destinationAccountID)
	if err != nil {
		entry.Errorf("Failed to load destinationAccount on transfer creation with internal error: %v", err)
		return &transfer.Transfer{}, err
	}
	if destinationAccount.ID == 0 {
		entry.Warn("Failed to create account with nonexisting destinationAccount")
		return &transfer.Transfer{}, account.ErrAccountNotFound
	}

	newTransfer, err := transfer.NewTransfer(originAccount, destinationAccount, amount)
	if err != nil {
		return &transfer.Transfer{}, err
	}

	err = u.transferRepo.CreateTransferAndUpdateAccountsBalances(ctx, newTransfer)
	if err != nil {
		logrus.New().WithField("transfer", newTransfer).Errorf("Failed to create transaction with internal error: %v", err)
		return &transfer.Transfer{}, err
	}

	return newTransfer, nil
}

func (u TransferUsecase) ListTransfers(ctx context.Context, loggedAccountID int) ([]*transfer.Transfer, error) {
	entry := logrus.WithField("accountID", loggedAccountID)

	transfers, err := u.transferRepo.GetTransfersCantainingAccount(ctx, loggedAccountID)
	if err != nil {
		entry.Errorf("Failed to list transfer while trying to load stored transfers: %v", err)
		return []*transfer.Transfer{}, nil
	}

	return transfers, nil
}
