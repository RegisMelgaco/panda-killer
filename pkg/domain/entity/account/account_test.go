package account_test

import (
	"fmt"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/shared"
	"testing"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewAccount(t *testing.T) {
	cases := []struct {
		testName    string
		name        string
		balance     shared.Money
		cpf         string
		secret      string
		expected    *account.Account
		expectedErr error
	}{
		{
			testName: "valid cpf and name SHOULD retrieve account.",
			name:     "Carminha",
			balance:  1,
			cpf:      "12345678901",
			secret:   ";)",
			expected: &account.Account{
				ID:        0,
				Balance:   shared.Money(1),
				Name:      "Maria",
				CPF:       "12345678901",
				Secret:    ";)",
				CreatedAt: time.Now(),
			},
			expectedErr: nil,
		},
		{
			testName: fmt.Sprintf("invalid cpf and valid name SHOULD retrive error %v", account.ErrAccountCPFShouldHaveLength11),
			name:     "Jully",
			balance:  3,
			cpf:      "123",
			secret:   ";)",
			expected: nil,
			expectedErr: validation.Errors{
				"CPF": account.ErrAccountCPFShouldHaveLength11,
			},
		},
		{
			testName: fmt.Sprintf("valid cpf and empty name SHOULD retrive error %v", account.ErrAccountNameIsObligatory),
			name:     "",
			balance:  42,
			cpf:      "12345678901",
			secret:   ";)",
			expected: nil,
			expectedErr: validation.Errors{
				"Name": account.ErrAccountNameIsObligatory,
			},
		},
		{
			testName: fmt.Sprintf("invalid cpf and name SHOULD retrive errors %v and %v", account.ErrAccountNameIsObligatory, account.ErrAccountCPFShouldHaveLength11),
			name:     "",
			balance:  42,
			cpf:      "123",
			secret:   ";)",
			expected: nil,
			expectedErr: validation.Errors{
				"Name": account.ErrAccountNameIsObligatory,
				"CPF":  account.ErrAccountCPFShouldHaveLength11,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			t.Parallel()
			actual, err := account.CreateNewAccount(c.balance, c.name, c.cpf, c.secret)

			if actual != nil && c.expected != nil {
				// Check if actual actual.CreatedAt is time.Now
				accoutCreatedAtDelta := actual.CreatedAt.Sub(c.expected.CreatedAt)
				assert.LessOrEqual(t, accoutCreatedAtDelta, time.Second)

				actual.CreatedAt = c.expected.CreatedAt
			}

			assert.Equal(t, c.expected, actual)

			switch expectedErrors := c.expectedErr.(type) {
			case validation.Errors:
				for field, expectedErr := range expectedErrors {
					assert.Equal(
						t,
						expectedErr.Error(),
						err.(validation.Errors)[field].(validation.ErrorObject).Message(),
					)
				}
			default:
				assert.Equal(t, c.expectedErr, err)
			}
		})
	}
}
