package e2etest

import (
	"context"
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

func TestLogin(t *testing.T) {
	ctx := context.Background()
	postgres.RunMigrations()

	pgxConn, _ := postgres.OpenConnection()
	accountRepo := repository.NewAccountRepo(pgxConn)
	transferRepo := repository.NewTransferRepo(pgxConn)
	passAlgo := algorithms.PasswordHashingAlgorithmsImpl{}
	accountUsecase := usecase.NewAccountUsecase(accountRepo, passAlgo)
	sessionAlgo := algorithms.SessionTokenAlgorithmsImpl{}
	router := rest.CreateRouter(
		usecase.NewAccountUsecase(accountRepo, passAlgo),
		usecase.NewTransferUsecase(transferRepo, accountRepo),
		usecase.NewAuthUsecase(accountRepo, sessionAlgo, passAlgo),
	)
	ts := httptest.NewServer(router)
	defer ts.Close()
	client := requests.Client{Host: ts.URL}

	t.Run("Login with success should receive response with session on headers", func(t *testing.T) {
		correctCPF := "12345678901"
		correctPassword := "pass"

		_, err := accountUsecase.CreateAccount(ctx, 0, "Joana", correctCPF, correctPassword)
		if err != nil {
			t.Errorf("Failed to create a test account: %v", err)
			t.FailNow()
		}

		resp, err := client.Login(rest.LoginRequest{
			CPF:      correctCPF,
			Password: correctPassword,
		})
		if err != nil {
			t.Errorf("Failed to request login: %v", err)
			t.FailNow()
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Didn't received expected status (OK): %v", resp.Status)
		}
		if len(resp.Header.Get("Authorization")) == 0 {
			t.Error("Authorization header should be set.")
		}
	})
	t.Run("Login without success should receive unauthorized", func(t *testing.T) {
		// t.Skip()
	})
}
