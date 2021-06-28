package e2etest

import (
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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	t.Run("Creating account successfully should persist account", func(t *testing.T) {
		postgres.RunMigrations()
		testStartTime := time.Now()

		testAccount := account.Account{
			Balance: 2,
			Name:    "Marcelinho",
			CPF:     "12345678901",
			Secret:  "s",
		}

		pgxConn, _ := postgres.OpenConnection()
		accountRepo := repository.NewAccountRepo(pgxConn)
		router := rest.CreateRouter(
			usecase.NewAccountUsecase(accountRepo),
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
		if respAccount.CreatedAt.Before(testStartTime) {
			t.Errorf("CreatedAt of account should be greatter or equals to test start time and not %v", respAccount.CreatedAt)
		}

		accounts, err := accountRepo.GetAccounts()
		if err != nil {
			t.Errorf("Failed to get stored accounts: %v", err)
			t.FailNow()
		}

		if len(accounts) != 1 {
			t.Errorf("Should exist one account")
			t.FailNow()
		}

		account := accounts[0]
		if account.ID != 1 {
			t.Errorf("Account id should be set and not %v", account.ID)
		}

		testAccount.ID = account.ID
		testAccount.CreatedAt = account.CreatedAt
		assert.ObjectsAreEqualValues(testAccount, account)
	})
	t.Run("Creating account with invalid account shouldn't persist account", func(t *testing.T) {
		postgres.RunMigrations()

		testAccount := account.Account{}

		pgxConn, _ := postgres.OpenConnection()
		accountRepo := repository.NewAccountRepo(pgxConn)
		router := rest.CreateRouter(
			usecase.NewAccountUsecase(accountRepo),
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
