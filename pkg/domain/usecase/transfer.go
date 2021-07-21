package usecase

import (
	"context"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/shared"
	"local/panda-killer/pkg/domain/entity/transfer"
	"sync"

	"github.com/sirupsen/logrus"
)

type TransferUsecase struct {
	transferRepo transfer.TransferRepo
	accountRepo  account.AccountRepo

	mux            *sync.Mutex
	accountMutexes map[account.AccountID]*sync.Mutex
}

func NewTransferUsecase(transferRepo transfer.TransferRepo, accountRepo account.AccountRepo) *TransferUsecase {
	return &TransferUsecase{
		transferRepo: transferRepo,
		accountRepo:  accountRepo,

		mux:            &sync.Mutex{},
		accountMutexes: make(map[account.AccountID]*sync.Mutex),
	}
}

func (u TransferUsecase) CreateTransfer(ctx context.Context, originAccountID, destinationAccountID account.AccountID, amount shared.Money) (*transfer.Transfer, error) {
	return u.handleCreateTransferParallelism(ctx, originAccountID, destinationAccountID, amount)
}

func (u TransferUsecase) handleCreateTransferParallelism(ctx context.Context, originAccountID, destinationAccountID account.AccountID, amount shared.Money) (*transfer.Transfer, error) {
	if originAccountID == destinationAccountID {
		return &transfer.Transfer{}, transfer.ErrTransferOriginAndDestinationNeedToBeDiffrent
	}

	u.mux.Lock()
	originMutex := u.getAccountMutexOrCreate(originAccountID)
	destinationMutex := u.getAccountMutexOrCreate(destinationAccountID)
	originMutex.Lock()
	destinationMutex.Lock()
	u.mux.Unlock()

	t, err := u.handleCreateTransferBussLogic(ctx, originAccountID, destinationAccountID, amount)

	originMutex.Unlock()
	destinationMutex.Unlock()

	u.mux.Lock()
	delete(u.accountMutexes, originAccountID)
	delete(u.accountMutexes, destinationAccountID)
	u.mux.Unlock()

	return t, err
}
func (u TransferUsecase) getAccountMutexOrCreate(accountID account.AccountID) *sync.Mutex {
	accountMutex, ok := u.accountMutexes[accountID]
	if !ok {
		accountMutex = &sync.Mutex{}
		u.accountMutexes[accountID] = accountMutex
	}
	return accountMutex
}

func (u TransferUsecase) handleCreateTransferBussLogic(ctx context.Context, originAccountID, destinationAccountID account.AccountID, amount shared.Money) (*transfer.Transfer, error) {
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

func (u TransferUsecase) ListTransfers(ctx context.Context, loggedAccountID account.AccountID) ([]transfer.Transfer, error) {
	entry := logrus.WithField("accountID", loggedAccountID)

	transfers, err := u.transferRepo.GetTransfersCantainingAccount(ctx, loggedAccountID)
	if err != nil {
		entry.Errorf("Failed to list transfer while trying to load stored transfers: %v", err)
		return []transfer.Transfer{}, nil
	}

	return transfers, nil
}
