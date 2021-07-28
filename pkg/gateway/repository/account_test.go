package repository_test

import (
	"context"
	"fmt"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/gateway/db/postgres/sqlc"
	"local/panda-killer/pkg/gateway/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	t.Parallel()

	_, _, pgPool := repository.CreateNewTestDBAndEnv(t.Name())
	q := sqlc.New(pgPool)
	repo := repository.NewAccountRepo(q)

	cases := []struct {
		testName                  string
		account                   *account.Account
		expectedErr               error
		idShouldBeGreaterThanZero bool
		setup                     func()
	}{
		{
			testName: "Valid account SHOULD change id to greater than 0",
			account: &account.Account{
				ID:        0,
				Balance:   42,
				Name:      "Mary",
				CPF:       "38333481868",
				Secret:    ";)",
				CreatedAt: time.Now(),
			},
			expectedErr:               nil,
			idShouldBeGreaterThanZero: true,
		},
		{
			testName: fmt.Sprintf("Repeated account cpf should retrive error %v and not update id", account.ErrAccountCPFShouldBeUnique),
			account: &account.Account{
				ID:        0,
				Balance:   42,
				Name:      "Katia",
				CPF:       "66828660382",
				Secret:    ";)",
				CreatedAt: time.Now(),
			},
			expectedErr:               account.ErrAccountCPFShouldBeUnique,
			idShouldBeGreaterThanZero: false,
			setup: func() {
				repo.CreateAccount(context.Background(), &account.Account{
					Name:      "Katia",
					CPF:       "66828660382",
					Secret:    ";)",
					CreatedAt: time.Now(),
				})
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			t.Parallel()

			c.setup()

			err := repo.CreateAccount(context.Background(), c.account)

			assert.Equal(t, c.expectedErr, err)

			if c.idShouldBeGreaterThanZero {
				assert.Greater(t, c.account.ID, account.AccountID(0))
			} else {
				assert.Equal(t, account.AccountID(0), c.account.ID)
			}
		})
	}
}

func TestGetAccountByCPF(t *testing.T) {
	t.Parallel()

	_, _, pgPool := repository.CreateNewTestDBAndEnv(t.Name())
	q := sqlc.New(pgPool)
	repo := repository.NewAccountRepo(q)

	cases := []struct {
		testName string

		expectedAccount *account.Account
		otherAccounts   []account.Account

		cpf         string
		expectedErr error
	}{
		{
			testName: "Get with existing cpf SHOULD get correct account",
			cpf:      "49008868580",
			expectedAccount: &account.Account{
				Name:      "Maria Jose",
				Balance:   99,
				CPF:       "49008868580",
				Secret:    ";)",
				CreatedAt: time.Now(),
			},
			otherAccounts: []account.Account{
				{
					Name:      "Joana",
					Balance:   36,
					CPF:       "49008868581",
					Secret:    ";)",
					CreatedAt: time.Now(),
				},
			},
		},
		{
			testName:        fmt.Sprintf("Get with nonexisting cpf SHOULD get error %v", account.ErrAccountNotFound),
			expectedAccount: nil,
			otherAccounts: []account.Account{
				{
					Name:      "Joana",
					Balance:   36,
					CPF:       "49008868573",
					Secret:    ";)",
					CreatedAt: time.Now(),
				},
			},
			cpf:         "0",
			expectedErr: account.ErrAccountNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			t.Parallel()

			accountsToBeStored := c.otherAccounts
			if c.expectedAccount != nil {
				accountsToBeStored = append(accountsToBeStored, *c.expectedAccount)
			}
			for _, a := range accountsToBeStored {
				repo.CreateAccount(context.Background(), &a)
			}

			actual, err := repo.GetAccountByCPF(context.Background(), c.cpf)

			if actual != nil {
				assert.Greater(t, actual.ID, 0)
				c.expectedAccount.ID = actual.ID

				createdAtDelta := c.expectedAccount.CreatedAt.Sub(actual.CreatedAt)
				assert.LessOrEqual(t, createdAtDelta, time.Second)
				c.expectedAccount.CreatedAt = actual.CreatedAt
			}

			assert.Equal(t, c.expectedAccount, actual)
			assert.Equal(t, c.expectedErr, err)
		})
	}
}
