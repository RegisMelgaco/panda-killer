package auth

import "local/panda-killer/pkg/domain/entity/account"

type PasswordHashingAlgorithms interface {
	GenerateSecretFromPassword(string) (string, error)
	// first string argument is secret and second is the password
	CheckSecretAndPassword(string, string) error
}

type SessionTokenAlgorithms interface {
	GenerateSessionToken(*account.Account) (string, error)
}
