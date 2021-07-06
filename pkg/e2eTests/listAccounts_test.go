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

	"github.com/stretchr/testify/assert"
)

func TestListAccounts(t *testing.T) {
	t.Run("List Accounts successfully should return persisted accounts", func(t *testing.T) {
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

		testAccounts := []account.Account{{Name: "João", CPF: "60684316730", Secret: "s"}, {Name: "Maria", CPF: "47577807613", Secret: "s"}}
		for _, a := range testAccounts {
			accountRepo.CreateAccount(context.Background(), &a)
		}

		resp, _ := client.ListAccounts()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code was OK and not %v", resp.Status)
		}
		var reqAccounts []account.Account
		err := json.NewDecoder(resp.Body).Decode(&reqAccounts)
		if err != nil {
			t.Errorf("Response could not be parsed: %v", err)
			t.FailNow()
		}
		assert.ObjectsAreEqualValues(reqAccounts, testAccounts)
	})
}
