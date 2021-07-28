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
