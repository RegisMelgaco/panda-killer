package account

import "time"

type Account struct {
	ID                int
	Balance           int
	Name, CPF, Secret string
	CreatedAt         time.Time
}

func checkAccountFieldsValidity(a *Account) error {
	if len(a.Name) == 0 {
		return ErrAccountNameIsObligatory
	}
	if len(a.CPF) != 11 {
		return ErrAccountCPFShouldHaveLength11
	}

	return nil
}

func CreateNewAccount(balance int, name, cpf, secret string) (*Account, error) {
	a := &Account{
		Balance:   balance,
		Name:      name,
		CPF:       cpf,
		Secret:    string(secret),
		CreatedAt: time.Now(),
	}

	err := checkAccountFieldsValidity(a)
	if err != nil {
		return &Account{}, err
	}

	return a, nil
}

func NewAccount(id, balance int, name, cpf, secret string, createdAt time.Time) *Account {
	return &Account{
		ID:        id,
		Balance:   balance,
		Name:      name,
		CPF:       cpf,
		Secret:    secret,
		CreatedAt: createdAt,
	}
}
