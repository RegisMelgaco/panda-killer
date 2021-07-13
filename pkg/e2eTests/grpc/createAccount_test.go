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
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateAccountGRPC(t *testing.T) {
	ctx := context.Background()
	postgres.RunMigrations()

	pgxConn, _ := postgres.OpenConnection()
	passAlgo := algorithms.PasswordHashingAlgorithmsImpl{}
	accountRepo := repository.NewAccountRepo(pgxConn)
	accountUsecase := usecase.NewAccountUsecase(accountRepo, passAlgo)

	s := rpc.NewApi(accountUsecase)

	t.Run("Creating account successfully should persist account", func(t *testing.T) {
		testAccount := &gen.CreateAccountRequest{
			Balance:  2,
			Name:     "Marcelinho",
			Cpf:      "74763452827",
			Password: "s",
		}

		respAccount, err := s.CreateAccount(ctx, testAccount)

		if err != nil {
			t.Errorf("Failed request to create account: %v", err)
			t.FailNow()
		}

		if respAccount.Id < 1 {
			t.Errorf("Id not set on response: %v", respAccount)
		}

		persistedAccounts, err := accountRepo.GetAccounts(context.Background())
		if err != nil {
			t.Errorf("Failed to get stored accounts: %v", err)
			t.FailNow()
		}

		if len(persistedAccounts) == 0 {
			t.Errorf("Should exist at least one account")
			t.FailNow()
		}

		persistedAccount := persistedAccounts[len(persistedAccounts)-1]

		if testAccount.Balance != int32(persistedAccount.Balance) ||
			testAccount.Cpf != persistedAccount.CPF ||
			testAccount.Name != persistedAccount.Name {
			t.Errorf("Persisted data doesn't match with request data: request = %v, persisted = %v", persistedAccount, testAccount)
		}

		if err = passAlgo.CheckSecretAndPassword(persistedAccount.Secret, testAccount.Password); err != nil {
			t.Errorf("Persisted secret from account should match to it's password.")
		}
	})

	t.Run("Create account with repeated cpf should retrieve BAD REQUEST and error", func(t *testing.T) {
		repeatedCPF := "12345678901"

		err := accountRepo.CreateAccount(context.Background(), &account.Account{
			Name:      "Joe",
			CPF:       repeatedCPF,
			Secret:    ";)",
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Errorf("Failed to create test account: %v", err)
			t.FailNow()
		}

		_, err = s.CreateAccount(context.Background(), &gen.CreateAccountRequest{
			Name: "Joe",
			Cpf:  repeatedCPF,
		})

		respStatus, _ := status.FromError(err)

		if respStatus.Code() != codes.InvalidArgument {
			t.Errorf("Expected response status code to be %v and not %v", codes.InvalidArgument, respStatus.Code())
		}

		if respStatus.Message() != account.ErrAccountCPFShouldBeUnique.Error() {
			t.Errorf("Expected response status message to be '%v' and not '%v'", account.ErrAccountCPFShouldBeUnique.Error(), respStatus.Message())
		}
	})

	t.Run("Creating account with cpf with length diffrent from 11 should retrive error and BAD REQUEST", func(t *testing.T) {
		testAccount := &gen.CreateAccountRequest{
			Cpf:  "123",
			Name: "Joe",
		}

		_, err := s.CreateAccount(ctx, testAccount)

		respStatus, _ := status.FromError(err)
		if respStatus.Code() != codes.InvalidArgument {
			t.Errorf("Server should answer with %v: %v", codes.InvalidArgument, respStatus.Code())
		}

		if respStatus.Message() != account.ErrAccountCPFShouldHaveLength11.Error() {
			t.Errorf("Received message diffrent from expected: expected=%v actual=%v", account.ErrAccountCPFShouldHaveLength11.Error(), respStatus.Message())
		}
	})
	postgres.DownToMigrationZero()
}
