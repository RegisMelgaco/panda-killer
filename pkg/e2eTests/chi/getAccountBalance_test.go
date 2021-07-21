package e2etest

import (
	"context"
	"encoding/json"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/shared"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/e2eTests/chi/requests"
	"local/panda-killer/pkg/gateway/algorithms"
	"local/panda-killer/pkg/gateway/db/postgres"
	"local/panda-killer/pkg/gateway/repository"
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountBalance(t *testing.T) {
	postgres.RunMigrations()

	pgxConn, _ := postgres.OpenConnection()
	accountRepo := repository.NewAccountRepo(pgxConn)
	transferRepo := repository.NewTransferRepo(pgxConn)
	passAlgo := algorithms.PasswordHashingAlgorithmsImpl{}
	sessionAlgo := algorithms.SessionTokenAlgorithmsImpl{}
	router := rest.CreateRouter(
		usecase.NewAccountUsecase(accountRepo, passAlgo),
		usecase.NewTransferUsecase(transferRepo, accountRepo),
		usecase.NewAuthUsecase(accountRepo, sessionAlgo, passAlgo),
	)
	ts := httptest.NewServer(router)
	defer ts.Close()
	client := requests.Client{Host: ts.URL}

	t.Run("Get account balance with success should retrive it's balance", func(t *testing.T) {
		var expectedBalance shared.Money = 42
		testAccount := account.Account{Name: "João", CPF: "34222086827", Secret: "s", Balance: expectedBalance}
		err := accountRepo.CreateAccount(context.Background(), &testAccount)
		if err != nil {
			t.Errorf("Failed to create test account: %v", err)
			t.FailNow()
		}

		resp, _ := client.GetAccountBalance(testAccount.ID)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Response status should be OK and not %v", resp.Status)
		}

		var balanceContainer rest.AccountBalanceResponse
		err = json.NewDecoder(resp.Body).Decode(&balanceContainer)
		if err != nil {
			t.Errorf("Request body format not expected: %v", err)
			t.FailNow()
		}

		if expectedBalance != balanceContainer.Balance {
			t.Errorf("Actual balance (%v) is diffrent from expected (%v)", balanceContainer.Balance, expectedBalance)
		}
	})
	t.Run("Get account balance from nonexisting account should retrieve a 404", func(t *testing.T) {
		resp, _ := client.GetAccountBalance(424242)

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected header was NOT FOUND and not %v", resp.Status)
		}
	})

	postgres.DownToMigrationZero()
}
