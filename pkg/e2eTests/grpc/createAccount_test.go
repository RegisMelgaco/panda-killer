package e2etest

import (
	"context"
	"fmt"
	"local/panda-killer/cmd/config"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/gateway/algorithms"
	"local/panda-killer/pkg/gateway/db/postgres"
	"local/panda-killer/pkg/gateway/db/postgres/sqlc"
	"local/panda-killer/pkg/gateway/repository"
	"local/panda-killer/pkg/gateway/rpc"
	"local/panda-killer/pkg/gateway/rpc/gen"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/ory/dockertest/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateAccountGRPC(t *testing.T) {
	ctx := context.Background()
	var env config.EnvVariablesProvider = config.EnvVariablesProviderImpl{}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Errorf("Could not connect to docker: %s", err)
		t.Fail()
	}

	dbName := t.Name()
	resource, err := pool.Run("postgres", "13.3", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=" + dbName})
	if err != nil {
		t.Errorf("Could not start resource: %s", err)
		t.Fail()
	}

	fmt.Println(resource.GetPort("5432/tcp"))

	env = env.SetTestDBUrl(
		fmt.Sprintf(
			"postgres://postgres:postgres@localhost:%s/%s?user=postgres&password=secret&sslmode=disable",
			resource.GetPort("5432/tcp"),
			dbName,
		),
	)

	var pgConn *pgx.Conn
	if err = pool.Retry(func() error {
		pgConn, err = postgres.OpenConnection(env)
		if err != nil {
			return err
		}
		return pgConn.Ping(ctx)
	}); err != nil {
		t.Errorf("Could not connect to docker: %s", err)
	}
	defer pgConn.Close(ctx)

	defer func() {
		if err = pool.Purge(resource); err != nil {
			t.Errorf("Could not purge resource: %s", err)
		}
	}()

	postgres.RunMigrations(env)

	pgPool, _ := postgres.OpenConnectionPool(env)
	defer pgPool.Close()
	queries := sqlc.New(pgPool)

	accountRepo := repository.NewAccountRepo(queries)
	passAlgo := algorithms.PasswordHashingAlgorithmsImpl{}
	accountUsecase := usecase.NewAccountUsecase(accountRepo, passAlgo)

	s := &rpc.Api{
		AccountUsecase: accountUsecase,
	}

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

	postgres.DownToMigrationZero(env)
}
