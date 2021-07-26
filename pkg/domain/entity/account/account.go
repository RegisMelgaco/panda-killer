package account

import (
	"local/panda-killer/pkg/domain/entity/shared"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Account struct {
	ID                AccountID
	Balance           shared.Money
	Name, CPF, Secret string
	CreatedAt         time.Time
}

var (
	cpfValidation  = validation.Length(11, 11).Error(ErrAccountCPFShouldHaveLength11.Error())
	nameValidation = validation.Required.Error(ErrAccountNameIsObligatory.Error())
)

func CreateNewAccount(balance shared.Money, name, cpf, secret string) (*Account, error) {
	a := Account{
		Balance:   balance,
		Name:      name,
		CPF:       cpf,
		Secret:    secret,
		CreatedAt: time.Now(),
	}
	err := validation.ValidateStruct(
		&a,
		validation.Field(&a.Name, nameValidation),
		validation.Field(&a.CPF, cpfValidation),
	)
	return &a, err
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
