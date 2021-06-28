package usecase

import (
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

func (u AccountUsecase) GetAccounts() ([]*account.Account, error) {
	return u.repo.GetAccounts()
}

func (u AccountUsecase) CreateAccount(newAccount *account.Account) error {
	if len(newAccount.Name) == 0 {
		return account.ErrAccountNameIsObligatory
	}
	if len(newAccount.CPF) != 11 {
		return account.ErrAccountCPFShouldHaveLength11
	}
	newAccount.CreatedAt = time.Now()
	return u.repo.CreateAccount(newAccount)
}
