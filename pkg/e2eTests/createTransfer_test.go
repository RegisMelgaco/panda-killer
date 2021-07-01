package e2etest

import (
	"context"
	"encoding/json"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/transfer"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/e2eTests/requests"
	"local/panda-killer/pkg/gateway/db/postgres"
	"local/panda-killer/pkg/gateway/repository"
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTransfer(t *testing.T) {
	t.Run("Create transfer with success should update users balances with success", func(t *testing.T) {
		err := postgres.RunMigrations()
		if err != nil {
			t.Errorf("Failed to run migrations: %v", err)
			t.FailNow()
		}

		pgxConn, _ := postgres.OpenConnection()
		accountRepo := repository.NewAccountRepo(pgxConn)
		transferRepo := repository.NewTransferRepo(pgxConn)
		router := rest.CreateRouter(
			usecase.NewAccountUsecase(accountRepo),
			usecase.NewTransferUsecase(transferRepo, accountRepo),
		)
		ts := httptest.NewServer(router)
		defer ts.Close()
		client := requests.Client{Host: ts.URL}

		testAccount1 := account.Account{Balance: 0.1, Name: "Maria", CPF: "12345678901"}
		err = accountRepo.CreateAccount(context.Background(), &testAccount1)
		if err != nil {
			t.Errorf("Failed to create test account1: %v", err)
		}

		testAccount2 := account.Account{Balance: 0.2, Name: "Joana", CPF: "12345678901"}
		err = accountRepo.CreateAccount(context.Background(), &testAccount2)
		if err != nil {
			t.Errorf("Failed to create test account2: %v", err)
			t.FailNow()
		}

		testTransfer := transfer.Transfer{AccountOrigin: testAccount1.ID, AccountDestination: testAccount2.ID, Amount: 0.1}
		resp, _ := client.CreateTransfer(testTransfer)

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Transfer creation request response should be OK not %v", resp.Status)
		}

		var respBodyObj rest.CreateTransferResponse
		err = json.NewDecoder(resp.Body).Decode(&respBodyObj)
		if err != nil {
			t.Errorf("Failed to decode response body: %v", err)
			t.FailNow()
		}

		if respBodyObj.ID < 1 {
			t.Errorf("Response body with invalid id: %v", respBodyObj.ID)
		}

		a1, _ := accountRepo.GetAccount(context.Background(), testAccount1.ID)
		if a1.Balance != 0 {
			t.Errorf("Expected balance for account 1 was 0 and not %v", a1.Balance)
		}
		a2, _ := accountRepo.GetAccount(context.Background(), testAccount2.ID)
		if a2.Balance != 0.3 {
			t.Errorf("Expected balance for account 2 was 0.3 and not %v", a2.Balance)
		}
	})
	t.Run("Create transfer without insufficient balance should fail", func(t *testing.T) {
		//TODO Implement
	})
	t.Run("Create transfer with non existing account(s) should fail", func(t *testing.T) {
		//TODO Implement
	})
}
