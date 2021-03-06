package e2etest

import (
	"context"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/e2eTests/chi/requests"
	"local/panda-killer/pkg/gateway/algorithms"
	"local/panda-killer/pkg/gateway/db/postgres/sqlc"
	"local/panda-killer/pkg/gateway/repository"
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	env, pgConn, pgPool := repository.CreateNewTestDBAndEnv(t.Name())
	queries := sqlc.New(pgPool)

	accountRepo := repository.NewAccountRepo(queries)
	transferRepo := repository.NewTransferRepo(pgConn)
	passAlgo := algorithms.PasswordHashingAlgorithmsImpl{}
	accountUsecase := usecase.NewAccountUsecase(accountRepo, passAlgo)
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

	correctCPF := "53744698351"
	correctPassword := "pass"

	_, err := accountUsecase.CreateAccount(ctx, 0, "Joana", correctCPF, correctPassword)
	if err != nil {
		t.Errorf("Failed to create a test account: %v", err)
		t.FailNow()
	}

	t.Run("Login with success should receive response with session on headers", func(t *testing.T) {
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
	t.Run("Login with incorrect password should receive unauthorized", func(t *testing.T) {
		incorrectPassword := correctPassword + "123"
		resp, err := client.Login(rest.LoginRequest{
			CPF:      correctCPF,
			Password: incorrectPassword,
		})
		if err != nil {
			t.Errorf("Failed to request login: %v", err)
			t.FailNow()
		}

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Didn't received expected status (Unauthorized): %v", resp.Status)
		}
		if len(resp.Header.Get("Authorization")) != 0 {
			t.Error("Authorization header should be set.")
		}
	})
	t.Run("Login with incorrect cpf should receive unauthorized", func(t *testing.T) {
		incorrectCPF := correctCPF + "123"
		resp, err := client.Login(rest.LoginRequest{
			CPF:      incorrectCPF,
			Password: correctPassword,
		})
		if err != nil {
			t.Errorf("Failed to request login: %v", err)
			t.FailNow()
		}

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Didn't received expected status (Unauthorized): %v", resp.Status)
		}
		if len(resp.Header.Get("Authorization")) != 0 {
			t.Error("Authorization header should be set.")
		}
	})
}
