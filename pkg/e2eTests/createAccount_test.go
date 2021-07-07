package e2etest

import (
	"context"
	"encoding/json"
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

	t.Run("Creating account successfully should persist account", func(t *testing.T) {
		testAccount := rest.CreateAccountRequest{
			Balance:  2,
			Name:     "Marcelinho",
			CPF:      "74763452827",
			Password: "s",
		}

		resp, _ := client.CreateAccount(testAccount)

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Failed request to create account: %v", resp)
			t.FailNow()
		}

		var respAccount rest.CreatedAccountResponse
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

		if len(persistedAccounts) == 0 {
			t.Errorf("Should exist at least one account")
			t.FailNow()
		}

		persistedAccount := persistedAccounts[len(persistedAccounts)-1]

		if testAccount.Balance != persistedAccount.Balance ||
			testAccount.CPF != persistedAccount.CPF ||
			testAccount.Name != persistedAccount.Name {
			t.Errorf("Persisted data doesn't match with request data: request = %v, persisted = %v", persistedAccount, testAccount)
		}

		if err = passAlgo.CheckSecretAndPassword(persistedAccount.Secret, testAccount.Password); err != nil {
			t.Errorf("Persisted secret from account should match to it's password.")
		}
	})
	t.Run("Creating account with invalid account shouldn't persist account", func(t *testing.T) {
		testAccount := rest.CreateAccountRequest{}

		resp, _ := client.CreateAccount(testAccount)

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Server should answer with bad request: %v", resp)
		}
	})
}
