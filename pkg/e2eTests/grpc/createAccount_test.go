package e2etest

import (
	"context"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/gateway/algorithms"
	"local/panda-killer/pkg/gateway/db/postgres"
	"local/panda-killer/pkg/gateway/repository"
	"local/panda-killer/pkg/gateway/rpc"
	"local/panda-killer/pkg/gateway/rpc/gen"
	"testing"
)

func TestCreateAccountGRPC(t *testing.T) {
	ctx := context.TODO()
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
}
