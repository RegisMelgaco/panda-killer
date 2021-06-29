package usecase

import (
	"context"
	"local/panda-killer/pkg/domain/entity/account"
	"time"

	"github.com/sirupsen/logrus"
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
	accounts, err := u.repo.GetAccounts(ctx)
	if err != nil {
		logrus.Errorf("Get accounts failed with internal error: %v", err)
		return accounts, err
	}

	logrus.Info("Get accounts succeeded")

	return accounts, err
}

func (u AccountUsecase) CreateAccount(ctx context.Context, newAccount *account.Account) error {
	entry := logrus.WithField("account", newAccount)

	var err error
	if len(newAccount.Name) == 0 {
		err = account.ErrAccountNameIsObligatory
	}
	if len(newAccount.CPF) != 11 {
		err = account.ErrAccountCPFShouldHaveLength11
	}

	if err != nil {
		entry.Infof("Create account failed with domain error: %v", err)
		return err
	}

	newAccount.CreatedAt = time.Now()

	err = u.repo.CreateAccount(ctx, newAccount)
	if err != nil {
		entry.Infof("Create account failed with internal error: %v", err)
		return err
	}

	entry.Info("Created account with success")
	return nil
}

func (u AccountUsecase) GetBalance(ctx context.Context, accountID int) (float64, error) {
	entry := logrus.WithField("accountID", accountID)

	a, err := u.repo.GetAccount(ctx, accountID)
	if a == nil {
		entry.Infof("Get balance failed with domain error: %v", err)
		return 0, account.ErrAccountNotFound
	}
	if err != nil {
		entry.Errorf("Get balance failed with internal error: %v", err)
		return 0, err
	}

	return a.Balance, nil
}
