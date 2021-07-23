package account

import (
	"local/panda-killer/pkg/domain/entity/shared"
	"time"
)

type Account struct {
	ID                AccountID
	Balance           shared.Money
	Name, CPF, Secret string
	CreatedAt         time.Time
}

func CreateNewAccount(balance shared.Money, name, cpf, secret string) *Account {
	return &Account{
		Balance:   balance,
		Name:      name,
		CPF:       cpf,
		Secret:    secret,
		CreatedAt: time.Now(),
	}
}

func NewAccount(id AccountID, balance shared.Money, name, cpf, secret string, createdAt time.Time) *Account {
	return &Account{
		ID:        id,
		Balance:   balance,
		Name:      name,
		CPF:       cpf,
		Secret:    secret,
		CreatedAt: createdAt,
	}
}
