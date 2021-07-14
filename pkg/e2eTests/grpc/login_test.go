package e2etest

import (
	"context"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/gateway/algorithms"
	"local/panda-killer/pkg/gateway/db/postgres"
	"local/panda-killer/pkg/gateway/repository"
	"local/panda-killer/pkg/gateway/rpc"
	"local/panda-killer/pkg/gateway/rpc/gen"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLogin(t *testing.T) {
	ctx := context.Background()
	postgres.RunMigrations()

	pgxConn, _ := postgres.OpenConnection()
	passAlgo := algorithms.PasswordHashingAlgorithmsImpl{}
	sessionAlgo := algorithms.SessionTokenAlgorithmsImpl{}
	accountRepo := repository.NewAccountRepo(pgxConn)
	accountUsecase := usecase.NewAccountUsecase(accountRepo, passAlgo)
	authUsecase := usecase.NewAuthUsecase(accountRepo, sessionAlgo, passAlgo)

	s := &rpc.Api{AuthUsecase: authUsecase}

	correctCPF := "53744698351"
	correctPassword := "pass"

	_, err := accountUsecase.CreateAccount(ctx, 0, "Joana", correctCPF, correctPassword)
	if err != nil {
		t.Errorf("Failed to create a test account: %v", err)
		t.Fail()
	}

	t.Run("Login with success should receive response with token", func(t *testing.T) {
		resp, err := s.Login(ctx, &gen.LoginRequest{
			Cpf:      correctCPF,
			Password: correctPassword,
		})

		respStatus, _ := status.FromError(err)
		if respStatus.Code() != codes.OK {
			t.Errorf("Didn't received expected status (%v): %v", codes.OK, respStatus.Code())
		}
		if len(resp.Token) == 0 {
			t.Error("Authorization header should be set.")
		}
	})
	t.Run("Login with incorrect password should receive unauthorized", func(t *testing.T) {
		incorrectPassword := correctPassword + "123"
		_, err := s.Login(ctx, &gen.LoginRequest{
			Cpf:      correctCPF,
			Password: incorrectPassword,
		})

		respStatus, _ := status.FromError(err)
		if respStatus.Code() != codes.Unauthenticated {
			t.Errorf("Didn't received expected status (%v): %v", codes.Unauthenticated, respStatus.Code())
		}
	})
	t.Run("Login with incorrect cpf should receive unauthorized", func(t *testing.T) {
		incorrectCPF := correctCPF + "123"
		_, err := s.Login(ctx, &gen.LoginRequest{
			Cpf:      incorrectCPF,
			Password: correctPassword,
		})

		respStatus, _ := status.FromError(err)
		if respStatus.Code() != codes.Unauthenticated {
			t.Errorf("Didn't received expected status (%v): %v", codes.Unauthenticated, respStatus.Code())
		}
	})

	postgres.DownToMigrationZero()
}
