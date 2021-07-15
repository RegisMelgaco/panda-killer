package usecase

import (
	"context"
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

func (u AuthUsecase) Login(ctx context.Context, cpf, password string) (authorizationHeader string, err error) {
	entry := logrus.WithFields(logrus.Fields{
		"cpf": cpf,
	})

	userAccount, err := u.accountRepo.GetAccountByCPF(ctx, cpf)
	if err != nil {
		return "", err
	}

	if userAccount.ID < 1 {
		return "", auth.ErrInvalidCredentials
	}

	err = u.passAlgo.CheckSecretAndPassword(userAccount.Secret, password)
	if err != nil {
		return "", auth.ErrInvalidCredentials
	}

	authorizationHeader, err = u.sessionAlgo.GenerateSessionToken(userAccount)
	if err != nil {
		entry.Errorf("Failed to login while creating a session token: %v", err)
		return "", err
	}
	return
}

func (u AuthUsecase) AddClaimsToContext(ctx context.Context, token string) (context.Context, error) {
	claims, err := u.sessionAlgo.GetClaims(token)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, auth.SessionContextKey, claims), nil
}
