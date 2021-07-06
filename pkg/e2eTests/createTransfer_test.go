package e2etest

import (
	"context"
	"encoding/json"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/transfer"
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

func TestCreateTransfer(t *testing.T) {
	postgres.RunMigrations()

	pgxConn, _ := postgres.OpenConnection()
	accountRepo := repository.NewAccountRepo(pgxConn)
	transferRepo := repository.NewTransferRepo(pgxConn)
	passAlgo := algorithms.PasswordHashingAlgorithmsImpl{}
	sessionAlgo := algorithms.SessionTokenAlgorithmsImpl{}
	accountUsecase := usecase.NewAccountUsecase(accountRepo, passAlgo)
	router := rest.CreateRouter(
		usecase.NewAccountUsecase(accountRepo, passAlgo),
		usecase.NewTransferUsecase(transferRepo, accountRepo),
		usecase.NewAuthUsecase(accountRepo, sessionAlgo, passAlgo),
	)
	ts := httptest.NewServer(router)
	defer ts.Close()
	client := requests.Client{Host: ts.URL}

	passord := ";)"
	testUser1, err := accountUsecase.CreateAccount(context.Background(), 1, "João", "85385023495", passord)
	if err != nil {
		t.Errorf("Failed to create test user: %v", err)
		t.FailNow()
	}

	resp, err := client.Login(rest.LoginRequest{
		CPF:      testUser1.CPF,
		Password: passord,
	})
	if err != nil {
		t.Errorf("Failed to request login: %v", err)
		t.FailNow()
	}

	authorizationToken := resp.Header.Get("Authorization")

	t.Run("Create transfer with success should update users balances with success", func(t *testing.T) {
		testAccount1 := account.Account{Balance: 1, Name: "Maria", CPF: "36063067680", Secret: "s"}
		err := accountRepo.CreateAccount(context.Background(), &testAccount1)
		if err != nil {
			t.Errorf("Failed to create test account1: %v", err)
		}

		testAccount2 := account.Account{Balance: 2, Name: "Joana", CPF: "46834635203", Secret: "s"}
		err = accountRepo.CreateAccount(context.Background(), &testAccount2)
		if err != nil {
			t.Errorf("Failed to create test account2: %v", err)
			t.FailNow()
		}

		transferRequest := rest.CreateTransferRequest{OriginAccountID: testAccount1.ID, DestinationAccountID: testAccount2.ID, Amount: 1}
		resp, _ := client.CreateTransfer(authorizationToken, transferRequest)

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
		if a2.Balance != 3 {
			t.Errorf("Expected balance for account 2 was 3 and not %v", a2.Balance)
		}
	})
	t.Run("Create transfer without insufficient balance should fail", func(t *testing.T) {
		originalOriginAccountBalance := 1
		originalDestineAccountBalance := 0

		testAccount1 := account.Account{Balance: originalOriginAccountBalance, Name: "Maria", CPF: "06316417772", Secret: "s"}
		err := accountRepo.CreateAccount(context.Background(), &testAccount1)
		if err != nil {
			t.Errorf("Failed to create test account1: %v", err)
		}

		testAccount2 := account.Account{Balance: originalDestineAccountBalance, Name: "Joana", CPF: "70465273858", Secret: "s"}
		err = accountRepo.CreateAccount(context.Background(), &testAccount2)
		if err != nil {
			t.Errorf("Failed to create test account2: %v", err)
			t.FailNow()
		}

		transferRequest := rest.CreateTransferRequest{OriginAccountID: testAccount1.ID, DestinationAccountID: testAccount2.ID, Amount: originalOriginAccountBalance + 1}
		resp, _ := client.CreateTransfer(authorizationToken, transferRequest)

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Transfer creation request response should be BAD REQUEST not %v", resp.Status)
		}

		var transferResponse rest.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&transferResponse)
		if err != nil {
			t.Errorf("Failed to parse create request response: %v", err)
			t.FailNow()
		}

		if transferResponse.Message != transfer.ErrInsufficientFundsToMakeTransaction.Error() {
			t.Errorf(
				"Response message should be %v and not %v",
				transfer.ErrInsufficientFundsToMakeTransaction.Error(),
				transferResponse.Message,
			)
		}

		updatedOriginAccount, err := accountRepo.GetAccount(context.Background(), testAccount1.ID)
		if err != nil {
			t.Errorf("Failed to retrieve updatedOriginAccount: %v", err)
		}
		updatedDestinationAccount, err := accountRepo.GetAccount(context.Background(), testAccount2.ID)
		if err != nil {
			t.Errorf("Failed to retrieve updatedDestinationAccount: %v", err)
		}

		if updatedOriginAccount.Balance != originalOriginAccountBalance {
			t.Errorf("It was expected to origin account not to change")
		}
		if updatedDestinationAccount.Balance != originalDestineAccountBalance {
			t.Errorf("It was expected to destination account not to change")
		}
	})
	t.Run("Create transfer with non existing account(s) should fail", func(t *testing.T) {
		transferRequest := rest.CreateTransferRequest{OriginAccountID: 132, DestinationAccountID: 13212, Amount: 1}
		resp, _ := client.CreateTransfer(authorizationToken, transferRequest)

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Transfer creation request response should be BAD REQUEST not %v", resp.Status)
		}

		var transferResponse rest.ErrorResponse
		err := json.NewDecoder(resp.Body).Decode(&transferResponse)
		if err != nil {
			t.Errorf("Failed to parse create request response: %v", err)
			t.FailNow()
		}

		if transferResponse.Message != account.ErrAccountNotFound.Error() {
			t.Errorf(
				"Response message should be %v and not %v",
				account.ErrAccountNotFound.Error(),
				transferResponse.Message,
			)
		}
	})
	t.Run("Create transfer with value lesser than zero should fail", func(t *testing.T) {
		testAccount1 := account.Account{Name: "Maria", CPF: "34414381401", Secret: "s"}
		err := accountRepo.CreateAccount(context.Background(), &testAccount1)
		if err != nil {
			t.Errorf("Failed to create test account1: %v", err)
		}

		testAccount2 := account.Account{Name: "Joana", CPF: "42462747478", Secret: "s"}
		err = accountRepo.CreateAccount(context.Background(), &testAccount2)
		if err != nil {
			t.Errorf("Failed to create test account2: %v", err)
			t.FailNow()
		}

		resp, err := client.CreateTransfer(authorizationToken, rest.CreateTransferRequest{
			OriginAccountID:      testAccount1.ID,
			DestinationAccountID: testAccount2.ID,
			Amount:               0,
		})
		if err != nil {
			t.Errorf("Failed to make request: %v", err)
			t.FailNow()
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Response status should be BAD REQUEST and not %v", resp.Status)
		}

		var respBody rest.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			t.Errorf("Failed to parse response body: %v", err)
			t.FailNow()
		}

		if respBody.Message != transfer.ErrTransferAmountShouldBeGreatterThanZero.Error() {
			t.Errorf("Expected message in response was '%v' and not '%v'", transfer.ErrTransferAmountShouldBeGreatterThanZero.Error(), respBody.Message)
		}
	})
	t.Run("Should not be possible to transfer to your self", func(t *testing.T) {
		testAccount := account.Account{Name: "Maria", CPF: "06268075730", Secret: "s"}
		err := accountRepo.CreateAccount(context.Background(), &testAccount)
		if err != nil {
			t.Errorf("Failed to create test account1: %v", err)
		}

		resp, err := client.CreateTransfer(authorizationToken, rest.CreateTransferRequest{
			OriginAccountID:      testAccount.ID,
			DestinationAccountID: testAccount.ID,
			Amount:               42,
		})
		if err != nil {
			t.Errorf("Failed to make request: %v", err)
			t.FailNow()
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Response status should be BAD REQUEST and not %v", resp.Status)
		}

		var respBody rest.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			t.Errorf("Failed to parse response body: %v", err)
			t.FailNow()
		}

		if respBody.Message != transfer.ErrTransferOriginAndDestinationNeedToBeDiffrent.Error() {
			t.Errorf("Expected message in response was '%v' and not '%v'", transfer.ErrTransferOriginAndDestinationNeedToBeDiffrent.Error(), respBody.Message)
		}
	})
}
