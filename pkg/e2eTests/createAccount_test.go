package e2etest

import (
	"context"
	"encoding/json"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/e2eTests/requests"
	"local/panda-killer/pkg/gateway/algorithms"
	"local/panda-killer/pkg/gateway/db/postgres"
	"local/panda-killer/pkg/gateway/repository"
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	t.Run("Creating account successfully should persist account", func(t *testing.T) {
		postgres.RunMigrations()

		testAccount := rest.CreateAccountRequest{
			Balance:  2,
			Name:     "Marcelinho",
			CPF:      "12345678901",
			Password: "s",
		}

		pgxConn, _ := postgres.OpenConnection()
		accountRepo := repository.NewAccountRepo(pgxConn)
		transferRepo := repository.NewTransferRepo(pgxConn)
		securityAlgo := algorithms.AccountSecurityAlgorithmsImpl{}
		router := rest.CreateRouter(
			usecase.NewAccountUsecase(accountRepo, securityAlgo),
			usecase.NewTransferUsecase(transferRepo, accountRepo),
		)
		ts := httptest.NewServer(router)
		defer ts.Close()
		client := requests.Client{Host: ts.URL}

		resp, _ := client.CreateAccount(testAccount)

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Failed request to create account: %v", resp)
			t.FailNow()
		}

		var respAccount account.Account
		err := json.NewDecoder(resp.Body).Decode(&respAccount)
		if err != nil {
			t.Errorf("Invalid response format: %v", err)
			t.FailNow()
		}

		if respAccount.ID < 1 {
			t.Errorf("Id not set on response: %v", respAccount)
		}

		persistedAccounts, err := accountRepo.GetAccounts(context.Background())
		if err != nil {
			t.Errorf("Failed to get stored accounts: %v", err)
			t.FailNow()
		}

		if len(persistedAccounts) > 1 {
			t.Errorf("Should exist at least one account")
			t.FailNow()
		}

		persistedAccount := persistedAccounts[len(persistedAccounts)-1]

		if testAccount.Balance != persistedAccount.Balance ||
			testAccount.CPF != persistedAccount.CPF ||
			testAccount.Name != persistedAccount.Name {
			t.Errorf("Persisted data doesn't match with request data: request = %v, persisted = %v", persistedAccount, testAccount)
		}

		if err = securityAlgo.CheckSecretAndPassword(persistedAccount.Secret, testAccount.Password); err != nil {
			t.Errorf("Persisted secret from account should match to it's password.")
		}
	})
	t.Run("Creating account with invalid account shouldn't persist account", func(t *testing.T) {
		postgres.RunMigrations()

		testAccount := rest.CreateAccountRequest{}

		pgxConn, _ := postgres.OpenConnection()
		accountRepo := repository.NewAccountRepo(pgxConn)
		transferRepo := repository.NewTransferRepo(pgxConn)
		router := rest.CreateRouter(
			usecase.NewAccountUsecase(accountRepo, algorithms.AccountSecurityAlgorithmsImpl{}),
			usecase.NewTransferUsecase(transferRepo, accountRepo),
		)
		ts := httptest.NewServer(router)
		defer ts.Close()
		client := requests.Client{Host: ts.URL}

		resp, _ := client.CreateAccount(testAccount)

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Server should answer with bad request: %v", resp)
		}
	})
}
