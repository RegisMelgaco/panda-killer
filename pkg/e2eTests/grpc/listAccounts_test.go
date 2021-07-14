package e2etest

import (
	"context"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/gateway/algorithms"
	"local/panda-killer/pkg/gateway/db/postgres"
	"local/panda-killer/pkg/gateway/repository"
	"local/panda-killer/pkg/gateway/rpc"
	"local/panda-killer/pkg/gateway/rpc/gen"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestListAccounts(t *testing.T) {
	ctx := context.Background()
	postgres.RunMigrations()

	pgxConn, _ := postgres.OpenConnection()
	passAlgo := algorithms.PasswordHashingAlgorithmsImpl{}
	accountRepo := repository.NewAccountRepo(pgxConn)
	accountUsecase := usecase.NewAccountUsecase(accountRepo, passAlgo)

	s := &rpc.Api{AccountUsecase: accountUsecase}

	t.Run("List Accounts successfully should return persisted accounts", func(t *testing.T) {
		testAccounts := []account.Account{{Name: "João", CPF: "60684316730", Secret: "s"}, {Name: "Maria", CPF: "47577807613", Secret: "s"}}
		for i, a := range testAccounts {
			accountRepo.CreateAccount(context.Background(), &a)
			testAccounts[i] = a
		}

		accountList, err := s.ListAccounts(ctx, &emptypb.Empty{})

		reqStatus, _ := status.FromError(err)
		if reqStatus.Code() != codes.OK {
			t.Errorf("Expected status code was %v and not %v", codes.OK, reqStatus.Code())
		}

		var testAccountsAsRequest gen.GetAccountListResponse
		for _, a := range testAccounts {
			testAccountsAsRequest.Accounts = append(testAccountsAsRequest.Accounts, &gen.GetAccountResponse{
				Id:   int32(a.ID),
				Name: a.Name,
				Cpf:  a.CPF,
			})
		}
		if !reflect.DeepEqual(accountList.Accounts, testAccountsAsRequest.Accounts) {
			t.Errorf("Expected accountList and testAccountsAsRequest to be equals: reqAccounts=%v testAccountsAsRequest=%v", accountList, testAccountsAsRequest)
		}
	})

	postgres.DownToMigrationZero()
}
