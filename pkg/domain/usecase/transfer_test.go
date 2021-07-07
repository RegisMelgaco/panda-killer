package usecase_test

import (
	"context"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/transfer"
	"local/panda-killer/pkg/domain/usecase"
	"testing"
	"time"
)

type mockedAccountRepo struct{}

func (m mockedAccountRepo) CreateAccount(context.Context, *account.Account) error {
	return nil
}
func (m mockedAccountRepo) GetAccounts(context.Context) ([]*account.Account, error) {
	return make([]*account.Account, 0), nil
}
func (m mockedAccountRepo) GetAccount(ctx context.Context, id int) (*account.Account, error) {
	return &account.Account{ID: id, Balance: 1000000}, nil
}
func (m mockedAccountRepo) GetAccountByCPF(context.Context, string) (*account.Account, error) {
	return &account.Account{}, nil
}

type count struct {
	num int
}

type mockedTransferRepo struct {
	count    *count
	waitChan chan bool
}

func (m mockedTransferRepo) CreateTransferAndUpdateAccountsBalances(ctx context.Context, t *transfer.Transfer) error {
	m.count.num += 1
	<-m.waitChan
	return nil
}
func (m mockedTransferRepo) GetTransfersCantainingAccount(ctx context.Context, accountID int) ([]*transfer.Transfer, error) {
	return []*transfer.Transfer{}, nil
}

func TestHandleCreateTransferParallelism(t *testing.T) {
	t.Run("Transactions with a shared account should wait each other", func(t *testing.T) {
		ctx := context.Background()
		account1 := 1
		account2 := 2
		account3 := 3

		transferRepo := mockedTransferRepo{count: &count{}, waitChan: make(chan bool), informStartWaiting: make(chan bool, 2)}

		transactionUsecase := usecase.NewTransferUsecase(transferRepo, mockedAccountRepo{})
		go transactionUsecase.CreateTransfer(ctx, account1, account2, 42)
		go transactionUsecase.CreateTransfer(ctx, account1, account3, 42)

		// Wait transaction go routines go sleep in final state
		time.Sleep(time.Millisecond * 10)

		if transferRepo.count.num == 2 {
			t.Errorf("It was expected to only one transfer to trigger CreateTransferAndUpdateAccountsBalances.")
		}
	})

	// Example: Transaction 1 = {Amount: 42, OriginAccount: 1, DestinationAccount: 2} and Transaction 2 = {Amount: 42, OriginAccount: 2, DestinationAccount: 1}
	t.Run("Transactions with same accounts and swapped possitions should not cause a dead lock", func(t *testing.T) {})

	t.Run("Transactions with diffrent accounts should not block each other", func(t *testing.T) {})
}
