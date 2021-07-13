package e2etest

import (
	"context"
	"encoding/json"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/e2eTests/chi/requests"
	"local/panda-killer/pkg/gateway/algorithms"
	"local/panda-killer/pkg/gateway/db/postgres"
	"local/panda-killer/pkg/gateway/repository"
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

		resp, err := client.CreateAccount(rest.CreateAccountRequest{
			Name: "Joe",
			CPF:  repeatedCPF,
		})
		if err != nil {
			t.Errorf("Failed to request account creation: %v", err)
			t.FailNow()
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected response status to be BAD REQUEST and not %v", resp.Status)
		}

		var errResp rest.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
			t.FailNow()
		}

		if errResp.Message != account.ErrAccountCPFShouldBeUnique.Error() {
			t.Errorf("Expected error response to be '%v' and not '%v'", account.ErrAccountCPFShouldBeUnique.Error(), errResp.Message)
		}
	})

	t.Run("Creating account with cpf with length diffrent from 11 should retrive error and BAD REQUEST", func(t *testing.T) {
		testAccount := rest.CreateAccountRequest{
			CPF:  "123",
			Name: "Joe",
		}

		resp, _ := client.CreateAccount(testAccount)

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Server should answer with bad request: %v", resp)
		}

		var errorResp rest.ErrorResponse
		err := json.NewDecoder(resp.Body).Decode(&errorResp)
		if err != nil {
			t.Errorf("Failed to parse response: %v", err)
			t.FailNow()
		}

		if errorResp.Message != account.ErrAccountCPFShouldHaveLength11.Error() {
			t.Errorf("Received message diffrent from expected: expected=%v actual=%v", account.ErrAccountCPFShouldHaveLength11.Error(), errorResp.Message)
		}
	})

	postgres.DownToMigrationZero()
}
