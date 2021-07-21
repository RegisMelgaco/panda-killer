package auth

import (
	"local/panda-killer/pkg/domain/entity/account"
)

type PasswordHashingAlgorithms interface {
	GenerateSecretFromPassword(secret string) (password string, err error)
	// first string argument is secret and second is the password
	CheckSecretAndPassword(secret string, password string) error
}

type SessionTokenAlgorithms interface {
	GenerateAuthorizationString(*account.Account) (string, error)
	GetClaims(token string) (*Claims, error)
}
