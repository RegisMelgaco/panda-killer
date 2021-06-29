package e2etest

import (
	"context"
	"encoding/json"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/e2eTests/requests"
	"local/panda-killer/pkg/gateway/db/postgres"
	"local/panda-killer/pkg/gateway/repository"
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountBalance(t *testing.T) {
	// caso de sucesso
	t.Run("Get account balance with success should retrive it's balance", func(t *testing.T) {
		postgres.RunMigrations()

		pgxConn, _ := postgres.OpenConnection()
		accountRepo := repository.NewAccountRepo(pgxConn)
		router := rest.CreateRouter(
			usecase.NewAccountUsecase(accountRepo),
		)
		ts := httptest.NewServer(router)
		defer ts.Close()
		client := requests.Client{Host: ts.URL}

		expectedBalance := 42.0
		testAccount := account.Account{Name: "João", CPF: "1235678901", Balance: expectedBalance}
		err := accountRepo.CreateAccount(context.Background(), &testAccount)
		if err != nil {
			t.Errorf("Failed to create test account: %v", err)
		}

		resp, _ := client.GetAccountBalance(testAccount.ID)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Response status should be OK and not %v", resp.Status)
		}

		var balanceContainer rest.BalanceContainer
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
		postgres.RunMigrations()

		pgxConn, _ := postgres.OpenConnection()
		accountRepo := repository.NewAccountRepo(pgxConn)
		router := rest.CreateRouter(
			usecase.NewAccountUsecase(accountRepo),
		)
		ts := httptest.NewServer(router)
		defer ts.Close()
		client := requests.Client{Host: ts.URL}

		resp, _ := client.GetAccountBalance(424242)

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected header was NOT FOUND and not %v", resp.Status)
		}
	})
}
