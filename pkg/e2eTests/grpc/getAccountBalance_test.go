package e2etest_test

import (
	"context"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/shared"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/gateway/algorithms"
	"local/panda-killer/pkg/gateway/db/postgres/sqlc"
	"local/panda-killer/pkg/gateway/repository"
	"local/panda-killer/pkg/gateway/rpc"
	"local/panda-killer/pkg/gateway/rpc/gen"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetAccountBalance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	_, _, pgPool := repository.CreateNewTestDBAndEnv(t.Name())

	queries := sqlc.New(pgPool)

	passAlgo := algorithms.PasswordHashingAlgorithmsImpl{}
	accountRepo := repository.NewAccountRepo(queries)
	accountUsecase := usecase.NewAccountUsecase(accountRepo, passAlgo)

	s := &rpc.Api{
		AccountUsecase: accountUsecase,
	}

	t.Run("Get account balance with success should retrive it's balance", func(t *testing.T) {

		var expectedBalance shared.Money = 42
		testAccount := account.Account{Name: "João", CPF: "34222086827", Secret: "s", Balance: expectedBalance}
		err := accountRepo.CreateAccount(context.Background(), &testAccount)
		if err != nil {
			t.Errorf("Failed to create test account: %v", err)
			t.FailNow()
		}

		resp, err := s.GetAccountBalance(ctx, &gen.GetAccountBalanceRequest{
			AccountId: int32(testAccount.ID),
		})

		respStatus, _ := status.FromError(err)

		if respStatus.Code() != codes.OK {
			t.Errorf("Response status should be OK and not %v", respStatus.Code())
		}

		if resp.Balance != int32(expectedBalance) {
			t.Errorf("Actual balance (%v) is diffrent from expected (%v)", resp.Balance, expectedBalance)
		}
	})
	t.Run("Get account balance from nonexisting account should retrieve a 404", func(t *testing.T) {
		_, err := s.GetAccountBalance(ctx, &gen.GetAccountBalanceRequest{AccountId: int32(424242)})

		reqStatus, _ := status.FromError(err)
		if reqStatus.Code() != codes.NotFound {
			t.Errorf("Expected header was %v and not %v", codes.NotFound, reqStatus.Code())
		}
	})
}
