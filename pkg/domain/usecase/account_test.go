package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/shared"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/gateway/algorithms"
	"testing"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func TestCreateAccount(t *testing.T) {
	workingRepo := account.AccountRepoMock{
		CreateAccountFunc: func(contextMoqParam context.Context, account *account.Account) error {
			account.ID = 1
			account.Secret = "secret"
			account.CreatedAt = time.Date(1, 1, 1, 1, 1, 1, 1, time.Local)
			return nil
		},
	}
	unexpectedError := errors.New("Oh my god :O! Soo unexpected...")

	cases := []struct {
		testName string

		repo *account.AccountRepoMock

		balance  shared.Money
		name     string
		cpf      string
		password string

		expected                    *account.Account
		expectedErr                 error
		expectCallRepoCreateAccount interface{}
	}{
		{
			testName: "Create account with valid account data and working repo SHOULD call repo and retrieve created account.",
			repo:     &workingRepo,
			balance:  3,
			name:     "Joana Patrícia Silva",
			cpf:      "08009347507",
			password: "1234567",
			expected: &account.Account{
				ID:        1,
				Balance:   3,
				Name:      "Joana Patrícia Silva",
				CPF:       "08009347507",
				Secret:    "secret",
				CreatedAt: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
			},
			expectedErr:                 nil,
			expectCallRepoCreateAccount: true,
		},
		{
			testName:                    "Create account with invalid account data and working repo SHOULD not call repo and retrieve a validation error.",
			repo:                        &workingRepo,
			balance:                     3,
			name:                        "",
			cpf:                         "08009347507",
			password:                    "1234567",
			expected:                    nil,
			expectedErr:                 validation.Errors{},
			expectCallRepoCreateAccount: false,
		},
		{
			testName: fmt.Sprintf("Create account with repeated cpf SHOULD retrive error %v", account.ErrAccountCPFShouldBeUnique),
			repo: &account.AccountRepoMock{
				CreateAccountFunc: func(contextMoqParam context.Context, a *account.Account) error {
					return account.ErrAccountCPFShouldBeUnique
				},
			},
			balance:                     42,
			name:                        "Cristina",
			cpf:                         "12345678901",
			password:                    ";)",
			expected:                    nil,
			expectedErr:                 account.ErrAccountCPFShouldBeUnique,
			expectCallRepoCreateAccount: true,
		},
		{
			testName: "Create account with unexpected error from repo SHOULD retrive error",
			repo: &account.AccountRepoMock{
				CreateAccountFunc: func(contextMoqParam context.Context, account *account.Account) error {
					return unexpectedError
				},
			},
			balance:                     13,
			name:                        "Clara",
			cpf:                         "12345678901",
			password:                    ";)",
			expected:                    nil,
			expectedErr:                 unexpectedError,
			expectCallRepoCreateAccount: true,
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			t.Parallel()

			uc := usecase.NewAccountUsecase(c.repo, algorithms.PasswordHashingAlgorithmsImpl{})
			actual, actualErr := uc.CreateAccount(context.Background(), c.balance, c.name, c.cpf, c.password)

			if actual != nil {
				assert.Equal(t, *c.expected, *actual)
			} else {
				assert.Equal(t, c.expected, actual)
			}
			assert.IsType(t, c.expectedErr, actualErr)

			calledRepo := len(c.repo.CreateAccountCalls()) > 0
			assert.Equal(t, c.expectCallRepoCreateAccount, calledRepo, "It was expected to call repo? %v", calledRepo)
		})
	}
}

func TestGetBalance(t *testing.T) {
	// unexpectedError := errors.New("Oh my god :O! Soo unexpected...")
	cases := []struct {
		testName   string
		mockedRepo *account.AccountRepoMock

		accountID account.AccountID

		expected    shared.Money
		expectedErr error
	}{
		{
			testName: "Get balance from exiting account with working repo SHOULD retrive balance",
			mockedRepo: &account.AccountRepoMock{
				GetAccountFunc: func(contextMoqParam context.Context, accountID account.AccountID) (*account.Account, error) {
					return &account.Account{
						Balance: 42,
					}, nil
				},
			},
			accountID:   1,
			expected:    shared.Money(42),
			expectedErr: nil,
		},
		// {
		// 	testName: fmt.Sprintf("Get balance from nonexisting account and working repo SHOULD retrive error %v", account.ErrAccountNotFound),
		// },
		// {
		// 	testName: "Get balance with repo with unexpected error SHOULD retrive error unexpected error",
		// },
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			t.Parallel()

			uc := usecase.NewAccountUsecase(
				c.mockedRepo, algorithms.PasswordHashingAlgorithmsImpl{},
			)

			actual, err := uc.GetBalance(context.Background(), c.accountID)

			assert.Equal(t, c.expected, actual)
			assert.Equal(t, c.expectedErr, err)
		})
	}
}
