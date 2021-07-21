package usecase_test

import (
	"context"
	"errors"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/transfer"
	"local/panda-killer/pkg/domain/usecase"
	"sync/atomic"
	"testing"
	"time"
)

type mockedAccountRepo struct{}

func (m mockedAccountRepo) CreateAccount(context.Context, *account.Account) error {
	return nil
}
func (m mockedAccountRepo) GetAccounts(context.Context) ([]account.Account, error) {
	return make([]account.Account, 0), nil
}
func (m mockedAccountRepo) GetAccount(ctx context.Context, id account.AccountID) (*account.Account, error) {
	return &account.Account{ID: account.AccountID(id), Balance: 1000000}, nil
}
func (m mockedAccountRepo) GetAccountByCPF(context.Context, string) (*account.Account, error) {
	return &account.Account{}, nil
}

type count struct {
	num int32
}

type mockedTransferRepo struct {
	count    *count
	waitChan chan bool
}

func (m mockedTransferRepo) CreateTransferAndUpdateAccountsBalances(ctx context.Context, t *transfer.Transfer) error {
	atomic.AddInt32(&m.count.num, 1)
	<-m.waitChan
	return nil
}
func (m mockedTransferRepo) GetTransfersCantainingAccount(ctx context.Context, accountID account.AccountID) ([]transfer.Transfer, error) {
	return []transfer.Transfer{}, nil
}

func TestHandleCreateTransferParallelism(t *testing.T) {
	t.Run("Transactions with a shared account should wait each other", func(t *testing.T) {
		ctx := context.Background()
		var account1 account.AccountID = 1
		var account2 account.AccountID = 2
		var account3 account.AccountID = 3

		transferRepo := mockedTransferRepo{count: &count{}, waitChan: make(chan bool)}

		transactionUsecase := usecase.NewTransferUsecase(transferRepo, mockedAccountRepo{})
		go transactionUsecase.CreateTransfer(ctx, account1, account2, 42)
		go transactionUsecase.CreateTransfer(ctx, account1, account3, 42)

		// Wait transaction go routines go sleep in final state
		time.Sleep(10 * time.Millisecond)

		if transferRepo.count.num != 1 {
			t.Errorf("It was expected to only one transfer to trigger CreateTransferAndUpdateAccountsBalances.")
		}

		transferRepo.waitChan <- true
		time.Sleep(10 * time.Millisecond)
		if transferRepo.count.num != 2 {
			t.Errorf("It was expected to only one transfer to trigger CreateTransferAndUpdateAccountsBalances.")
		}
	})

	// Example: Transaction 1 = {Amount: 42, OriginAccount: 1, DestinationAccount: 2} and Transaction 2 = {Amount: 42, OriginAccount: 2, DestinationAccount: 1}
	t.Run("Transactions with same accounts and swapped possitions should not cause a dead lock", func(t *testing.T) {
		ctx := context.Background()
		var account1 account.AccountID = 1
		var account2 account.AccountID = 2

		transferRepo := mockedTransferRepo{count: &count{}, waitChan: make(chan bool)}

		transactionUsecase := usecase.NewTransferUsecase(transferRepo, mockedAccountRepo{})
		go transactionUsecase.CreateTransfer(ctx, account1, account2, 42)
		go transactionUsecase.CreateTransfer(ctx, account2, account1, 42)

		// Wait transaction go routines go sleep in final state
		time.Sleep(time.Millisecond * 10)
		if transferRepo.count.num != 1 {
			t.Errorf("It was expected to only one transfer to trigger CreateTransferAndUpdateAccountsBalances. %v", transferRepo.count.num)
		}

		// Free fist transaction go routine then wait to final state
		transferRepo.waitChan <- true
		time.Sleep(time.Millisecond * 10)
		if transferRepo.count.num != 2 {
			t.Errorf("It was expected to both transfers to trigger CreateTransferAndUpdateAccountsBalances. Transfers calling repo: %v", transferRepo.count.num)
		}
	})

	t.Run("Transactions with diffrent accounts should not block each other", func(t *testing.T) {
		ctx := context.Background()
		var account1 account.AccountID = 1
		var account2 account.AccountID = 2
		var account3 account.AccountID = 3
		var account4 account.AccountID = 4

		transferRepo := mockedTransferRepo{count: &count{}, waitChan: make(chan bool)}

		transactionUsecase := usecase.NewTransferUsecase(transferRepo, mockedAccountRepo{})
		go transactionUsecase.CreateTransfer(ctx, account1, account2, 42)
		go transactionUsecase.CreateTransfer(ctx, account3, account4, 42)

		// Wait transaction go routines go sleep in final state
		time.Sleep(time.Millisecond * 10)
		if transferRepo.count.num != 2 {
			t.Errorf("It was expected to only one transfer to trigger CreateTransferAndUpdateAccountsBalances.")
		}
	})

	t.Run("Transfer with it self should not cause a deadlock", func(t *testing.T) {
		ctx := context.Background()
		var account1 account.AccountID = 1

		transferRepo := mockedTransferRepo{count: &count{}, waitChan: make(chan bool)}

		transactionUsecase := usecase.NewTransferUsecase(transferRepo, mockedAccountRepo{})
		_, err := transactionUsecase.CreateTransfer(ctx, account1, account1, 42)
		if !errors.Is(err, transfer.ErrTransferOriginAndDestinationNeedToBeDiffrent) {
			t.Errorf("Transactions with destination account equals to origin account should return error %v", transfer.ErrTransferOriginAndDestinationNeedToBeDiffrent)
		}
	})
}
