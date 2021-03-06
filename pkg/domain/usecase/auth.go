package usecase

import (
	"context"
	"errors"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/auth"

	"github.com/sirupsen/logrus"
)

type AuthUsecase struct {
	accountRepo account.AccountRepo
	sessionAlgo auth.SessionTokenAlgorithms
	passAlgo    auth.PasswordHashingAlgorithms
}

func NewAuthUsecase(accountRepo account.AccountRepo, sessionAlgo auth.SessionTokenAlgorithms, passAlgo auth.PasswordHashingAlgorithms) *AuthUsecase {
	return &AuthUsecase{
		accountRepo: accountRepo,
		sessionAlgo: sessionAlgo,
		passAlgo:    passAlgo,
	}
}

func (u AuthUsecase) Login(ctx context.Context, cpf, password string) (authorization string, err error) {
	entry := logrus.WithFields(logrus.Fields{
		"cpf": cpf,
	})

	userAccount, err := u.accountRepo.GetAccountByCPF(ctx, cpf)
	if errors.Is(err, account.ErrAccountNotFound) {
		return "", auth.ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}

	err = u.passAlgo.CheckSecretAndPassword(userAccount.Secret, password)
	if err != nil {
		return "", auth.ErrInvalidCredentials
	}

	authorization, err = u.sessionAlgo.GenerateAuthorizationString(userAccount)
	if err != nil {
		entry.Errorf("Failed to login while creating a session token: %v", err)
		return "", err
	}
	return
}

func (u AuthUsecase) AddClaimsToContext(ctx context.Context, authorization string) (context.Context, error) {
	claims, err := u.sessionAlgo.GetClaims(authorization)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, auth.SessionContextKey, claims), nil
}
