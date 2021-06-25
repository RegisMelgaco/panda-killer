package usecase

import (
	"local/panda-killer/pkg/domain/entity/account"
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
	return u.repo.CreateAccount(newAccount)
}
