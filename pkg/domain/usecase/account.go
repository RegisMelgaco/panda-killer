package usecase

import (
	"context"
	"local/panda-killer/pkg/domain/entity/account"
	"time"
)

type AccountUsecase struct {
	repo account.AccountRepo
}

func NewAccountUsecase(accountRepo account.AccountRepo) *AccountUsecase {
	return &AccountUsecase{
		repo: accountRepo,
	}
}

func (u AccountUsecase) GetAccounts(ctx context.Context) ([]*account.Account, error) {
	return u.repo.GetAccounts(ctx)
}

func (u AccountUsecase) CreateAccount(ctx context.Context, newAccount *account.Account) error {
	if len(newAccount.Name) == 0 {
		return account.ErrAccountNameIsObligatory
	}
	if len(newAccount.CPF) != 11 {
		return account.ErrAccountCPFShouldHaveLength11
	}
	newAccount.CreatedAt = time.Now()
	return u.repo.CreateAccount(ctx, newAccount)
}

func (u AccountUsecase) GetBalance(ctx context.Context, accountID int) (float64, error) {
	return u.repo.GetAccountBalance(ctx, accountID)
}
