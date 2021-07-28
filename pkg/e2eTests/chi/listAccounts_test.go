package e2etest

import (
	"context"
	"encoding/json"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/e2eTests/chi/requests"
	"local/panda-killer/pkg/gateway/algorithms"
	"local/panda-killer/pkg/gateway/db/postgres/sqlc"
	"local/panda-killer/pkg/gateway/repository"
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestListAccounts(t *testing.T) {
	t.Parallel()

	env, pgConn, pgPool := repository.CreateNewTestDBAndEnv(t.Name())
	queries := sqlc.New(pgPool)

	accountRepo := repository.NewAccountRepo(queries)
	transferRepo := repository.NewTransferRepo(pgConn)
	passAlgo := algorithms.PasswordHashingAlgorithmsImpl{}
	sessionAlgo := algorithms.NewSessionTokenAlgorithms(env)
	router := rest.CreateRouter(
		env,
		usecase.NewAccountUsecase(accountRepo, passAlgo),
		usecase.NewTransferUsecase(transferRepo, accountRepo),
		usecase.NewAuthUsecase(accountRepo, sessionAlgo, passAlgo),
	)
	ts := httptest.NewServer(router)
	defer ts.Close()
	client := requests.Client{Host: ts.URL}

	t.Run("List Accounts successfully should return persisted accounts", func(t *testing.T) {
		testAccounts := []account.Account{{Name: "João", CPF: "60684316730", Secret: "s"}, {Name: "Maria", CPF: "47577807613", Secret: "s"}}
		for i, a := range testAccounts {
			accountRepo.CreateAccount(context.Background(), &a)
			testAccounts[i] = a
		}

		resp, _ := client.ListAccounts()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code was OK and not %v", resp.Status)
		}
		var reqAccounts []rest.GetAccountResponse
		err := json.NewDecoder(resp.Body).Decode(&reqAccounts)
		if err != nil {
			t.Errorf("Response could not be parsed: %v", err)
			t.FailNow()
		}

		var testAccountsAsRequest []rest.GetAccountResponse
		for _, a := range testAccounts {
			testAccountsAsRequest = append(testAccountsAsRequest, rest.GetAccountResponse{
				ID:   a.ID,
				Name: a.Name,
				CPF:  a.CPF,
			})
		}
		if !reflect.DeepEqual(reqAccounts, testAccountsAsRequest) {
			t.Errorf("Expected reqAccounts and testAccountsAsRequest to be equals: reqAccounts=%v testAccountsAsRequest=%v", reqAccounts, testAccountsAsRequest)
		}
	})
}
