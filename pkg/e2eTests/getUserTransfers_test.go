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
	"reflect"
	"testing"
)

func TestGetUserTransfers(t *testing.T) {
	ctx := context.Background()
	postgres.RunMigrations()

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

	t.Run("Get user transfers with success should return list of transfers where the user is part.", func(t *testing.T) {
		testUser1 := account.Account{Name: "João", CPF: "1234578901", Balance: 1}
		err := accountRepo.CreateAccount(ctx, &testUser1)
		if err != nil {
			t.Errorf("Failed to create test user: %v", err)
			t.FailNow()
		}
		testUser2 := account.Account{Name: "Malaquias", CPF: "1234578901"}
		err = accountRepo.CreateAccount(ctx, &testUser2)
		if err != nil {
			t.Errorf("Failed to create testUser2: %v", err)
			t.FailNow()
		}
		testUser3 := account.Account{Name: "Jorge", CPF: "1234578901"}
		err = accountRepo.CreateAccount(ctx, &testUser3)
		if err != nil {
			t.Errorf("Failed to create testUser3: %v", err)
			t.FailNow()
		}

		transfer1, err := transfer.NewTransfer(&testUser1, &testUser2, 1)
		if err != nil {
			t.Errorf("Failed to create transfer1: %v", err)
			t.FailNow()
		}
		err = transferRepo.CreateTransferAndUpdateAccountsBalances(
			ctx,
			transfer1,
		)
		if err != nil {
			t.Errorf("Failed to persist transfer1: %v", err)
			t.FailNow()
		}

		transfer2, err := transfer.NewTransfer(&testUser2, &testUser3, 1)
		if err != nil {
			t.Errorf("Failed to create transfer1: %v", err)
			t.FailNow()
		}
		err = transferRepo.CreateTransferAndUpdateAccountsBalances(
			ctx,
			transfer2,
		)
		if err != nil {
			t.Errorf("Failed to persist transfer2: %v", err)
			t.FailNow()
		}

		transfer3, err := transfer.NewTransfer(&testUser3, &testUser1, 1)
		if err != nil {
			t.Errorf("Failed to create transfer3: %v", err)
			t.FailNow()
		}
		err = transferRepo.CreateTransferAndUpdateAccountsBalances(
			ctx,
			transfer3,
		)
		if err != nil {
			t.Errorf("Failed to persist transfer3: %v", err)
			t.FailNow()
		}

		resp, err := client.ListTransfers(testUser1.ID)
		if err != nil {
			t.Errorf("Failed to request transfers list: %v", err)
			t.FailNow()
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected response status was OK and not %v", resp.Status)
		}

		var requestedTransfers []rest.GetTransferResponse
		err = json.NewDecoder(resp.Body).Decode(&requestedTransfers)
		if err != nil {
			t.Errorf("Failed to decode list transfers response")
			t.FailNow()
		}

		expectedResponse := []rest.GetTransferResponse{
			{
				ID:                   transfer1.ID,
				Amount:               transfer1.Amount,
				OriginAccountID:      transfer1.OriginAccount.ID,
				DestinationAccountID: transfer1.DestinationAccount.ID,
				CreatedAt:            requestedTransfers[0].CreatedAt,
			}, {
				ID:                   transfer3.ID,
				Amount:               transfer3.Amount,
				OriginAccountID:      transfer3.OriginAccount.ID,
				DestinationAccountID: transfer3.DestinationAccount.ID,
				CreatedAt:            requestedTransfers[1].CreatedAt,
			},
		}
		if !reflect.DeepEqual(expectedResponse, requestedTransfers) {
			t.Errorf("Expected request response body to be %v and not %v", expectedResponse, requestedTransfers)
		}
	})
}
