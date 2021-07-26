package usecase

import (
	"context"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/auth"
	"local/panda-killer/pkg/domain/entity/shared"

	"github.com/sirupsen/logrus"
)

type AccountUsecase struct {
	repo         account.AccountRepo
	securityAlgo auth.PasswordHashingAlgorithms
}

func NewAccountUsecase(accountRepo account.AccountRepo, securityAlgo auth.PasswordHashingAlgorithms) *AccountUsecase {
	return &AccountUsecase{
		repo:         accountRepo,
		securityAlgo: securityAlgo,
	}
}

func (u AccountUsecase) GetAccounts(ctx context.Context) ([]account.Account, error) {
	accounts, err := u.repo.GetAccounts(ctx)
	if err != nil {
		logrus.Errorf("Get accounts failed with internal error: %v", err)
		return accounts, err
	}

	logrus.Info("Get accounts succeeded")

	return accounts, err
}

func (u AccountUsecase) CreateAccount(ctx context.Context, balance shared.Money, name string, cpf string, password string) (*account.Account, error) {
	entry := logrus.WithFields(logrus.Fields{
		"balance": balance, "name": name, "cpf": cpf,
	})

	secret, err := u.securityAlgo.GenerateSecretFromPassword(password)
	if err != nil {
		entry.Errorf("Failed to create account while genereting secret: %v", err)
		return nil, err
	}

	newAccount, err := account.CreateNewAccount(balance, name, cpf, secret)
	if err != nil {
		entry.Infof("Create account failed with domain error: %v", err)
		return nil, err
	}

	err = u.repo.CreateAccount(ctx, newAccount)
	if err != nil {
		entry.Infof("Create account failed with internal error: %v", err)
		return nil, err
	}

	entry.Info("Created account with success")
	return newAccount, nil
}

func (u AccountUsecase) GetBalance(ctx context.Context, accountID account.AccountID) (shared.Money, error) {
	entry := logrus.WithField("accountID", accountID)

	a, err := u.repo.GetAccount(ctx, accountID)
	if a.ID == 0 {
		entry.Infof("Get balance failed with domain error: %v", err)
		return 0, account.ErrAccountNotFound
	}
	if err != nil {
		entry.Errorf("Get balance failed with internal error: %v", err)
		return 0, err
	}

	return a.Balance, nil
}
