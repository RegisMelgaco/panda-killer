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
		postgres.RunMigrations()

		pgxConn, _ := postgres.OpenConnection()
		accountRepo := repository.NewAccountRepo(pgxConn)
		router := rest.CreateRouter(
			usecase.NewAccountUsecase(accountRepo),
		)
		ts := httptest.NewServer(router)
		defer ts.Close()
		client := requests.Client{Host: ts.URL}

		testAccount1 := account.Account{Balance: 0.1, Name: "Maria", CPF: "12345678901"}
		err := accountRepo.CreateAccount(context.Background(), &testAccount1)
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

		if resp.StatusCode != http.StatusOK {
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
	})
	t.Run("Create transfer without insufficient balance should fail", func(t *testing.T) {
		//TODO Implement
	})
	t.Run("Create transfer with not existing account(s) should fail", func(t *testing.T) {
		//TODO Implement
	})
}
