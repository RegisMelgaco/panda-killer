package rpc_test

import (
	"context"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/shared"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/gateway/rpc"
	"local/panda-killer/pkg/gateway/rpc/gen"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testAccount := &account.Account{
		ID:        1,
		Balance:   101,
		Name:      "Valeria",
		CPF:       "1235678901",
		Secret:    ";)",
		CreatedAt: time.Now(),
	}

	cases := []struct {
		testName string

		api        rpc.Api
		accountReq *gen.CreateAccountRequest

		expected    *gen.CreateAccountResponse
		expectedErr error
	}{
		{
			testName: "Create account without error from usecase should retrieve account and status with ok",
			api: rpc.Api{
				AccountUsecase: &usecase.AccountUsecaseMock{
					CreateAccountFunc: func(ctx context.Context, balance shared.Money, name, cpf, password string) (*account.Account, error) {
						return testAccount, nil
					},
				},
			},
			accountReq: &gen.CreateAccountRequest{},
			expected: &gen.CreateAccountResponse{
				Id: int32(testAccount.ID),
			},
			expectedErr: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			t.Parallel()

			actual, err := c.api.CreateAccount(ctx, c.accountReq)

			assert.Equal(t, c.expected, actual)
			assert.Equal(t, c.expectedErr, err)
		})
	}
}
