package usecase_test

import (
	"context"
	"errors"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/gateway/algorithms"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccounts(t *testing.T) {
	cases := []struct {
		name        string
		mockedRepo  account.AccountRepo
		expected    []account.Account
		expectedErr error
	}{
		{
			name: "with repo mock without error should retrieve accounts from mock",
			mockedRepo: &account.AccountRepoMock{
				GetAccountsFunc: func(contextMoqParam context.Context) ([]account.Account, error) {
					return []account.Account{
						{
							ID:      1,
							Balance: 10,
							Name:    "Bruna",
							CPF:     "12345678912",
							Secret:  "78312798419",
						},
					}, nil
				},
			},
			expected: []account.Account{
				{
					ID:      1,
					Balance: 10,
					Name:    "Bruna",
					CPF:     "12345678912",
					Secret:  "78312798419",
				},
			},
			expectedErr: nil,
		},
		{
			name: "with repo mock with error should retrieve the error",
			mockedRepo: &account.AccountRepoMock{
				GetAccountsFunc: func(contextMoqParam context.Context) ([]account.Account, error) {
					return nil, errors.New("Internal error")
				},
			},
			expected:    nil,
			expectedErr: errors.New("Internal error"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			uc := usecase.NewAccountUsecase(
				c.mockedRepo, algorithms.PasswordHashingAlgorithmsImpl{},
			)

			accounts, err := uc.GetAccounts(context.Background())

			assert.Equal(t, c.expected, accounts)
			assert.Equal(t, c.expectedErr, err)
		})
	}
}
