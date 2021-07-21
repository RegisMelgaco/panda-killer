package e2etest

import (
	"context"
	"local/panda-killer/pkg/domain/entity/transfer"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/gateway/algorithms"
	"local/panda-killer/pkg/gateway/db/postgres"
	"local/panda-killer/pkg/gateway/db/postgres/sqlc"
	"local/panda-killer/pkg/gateway/repository"
	"local/panda-killer/pkg/gateway/rpc"
	"local/panda-killer/pkg/gateway/rpc/gen"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetUserTransfers(t *testing.T) {
	ctx := context.Background()

	postgres.RunMigrations()
	defer postgres.DownToMigrationZero()

	pgxConn, _ := postgres.OpenConnection()
	defer pgxConn.Close(context.Background())
	pgPool, _ := postgres.OpenConnectionPool()
	defer pgPool.Close()
	queries := sqlc.New(pgPool)

	passAlgo := algorithms.PasswordHashingAlgorithmsImpl{}
	sessionAlgo := algorithms.SessionTokenAlgorithmsImpl{}
	accountRepo := repository.NewAccountRepo(queries)
	transferRepo := repository.NewTransferRepo(pgxConn)
	accountUsecase := usecase.NewAccountUsecase(accountRepo, passAlgo)
	authUsecase := usecase.NewAuthUsecase(accountRepo, sessionAlgo, passAlgo)
	transferUsecase := usecase.NewTransferUsecase(transferRepo, accountRepo)

	api := &rpc.Api{
		AuthUsecase:     authUsecase,
		AccountUsecase:  accountUsecase,
		TransferUsecase: transferUsecase,
	}

	client, conn, err := StartServerAndGetClient(ctx, api)
	if err != nil {
		t.Errorf("Failed to start test server: %v", err)
	}
	defer conn.Close()

	t.Run("Get user transfers with success should return list of transfers where the user is part.", func(t *testing.T) {
		passord := ";)"
		testUser1, err := accountUsecase.CreateAccount(ctx, 1, "João", "33228058166", passord)
		if err != nil {
			t.Errorf("Failed to create test user: %v", err)
			t.FailNow()
		}
		testUser2, err := accountUsecase.CreateAccount(ctx, 0, "Malaquias", "94785942214", ";)")
		if err != nil {
			t.Errorf("Failed to create test user: %v", err)
			t.FailNow()
		}
		testUser3, err := accountUsecase.CreateAccount(ctx, 0, "Jorge", "03729912343", ";)")
		if err != nil {
			t.Errorf("Failed to create test user: %v", err)
			t.FailNow()
		}

		transfer1, err := transfer.NewTransfer(testUser1, testUser2, 1)
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

		transfer2, err := transfer.NewTransfer(testUser2, testUser3, 1)
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

		transfer3, err := transfer.NewTransfer(testUser3, testUser1, 1)
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

		loginResp, err := client.Login(ctx, &gen.LoginRequest{
			Cpf:      testUser1.CPF,
			Password: passord,
		})

		if _, ok := status.FromError(err); !ok {
			t.Errorf("Failed to request login: %v", err)
			t.FailNow()
		}

		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", loginResp.Token)

		resp, err := client.ListTransfers(ctx, &emptypb.Empty{})
		reqStatus, ok := status.FromError(err)
		if !ok {
			t.Errorf("Failed to request transfers list: %v", reqStatus.Details()...)
			t.FailNow()
		}

		if _, ok := status.FromError(err); !ok {
			t.Errorf("Expected response status was OK")
		}

		expectedResponse := []*gen.GetTranferResponse{
			{
				Id:                   int32(transfer1.ID),
				Amount:               int32(transfer1.Amount),
				OriginAccountId:      int32(transfer1.OriginAccount.ID),
				DestinationAccountId: int32(transfer1.DestinationAccount.ID),
				CreatedAt:            timestamppb.New(resp.Transfers[0].CreatedAt.AsTime()),
			}, {
				Id:                   int32(transfer3.ID),
				Amount:               int32(transfer3.Amount),
				OriginAccountId:      int32(transfer3.OriginAccount.ID),
				DestinationAccountId: int32(transfer3.DestinationAccount.ID),
				CreatedAt:            timestamppb.New(resp.Transfers[1].CreatedAt.AsTime()),
			},
		}
		if !reflect.DeepEqual(expectedResponse, resp.Transfers) {
			t.Errorf("Expected request response to be %v and not %v", expectedResponse, resp.Transfers)
		}
	})

	t.Run("Request transfers while unlogged should retrive unauthorized", func(t *testing.T) {
		ctx = metadata.AppendToOutgoingContext(context.Background(), "authorization", "invalid token")
		_, err := client.ListTransfers(ctx, &emptypb.Empty{})

		reqStatus, _ := status.FromError(err)
		if reqStatus.Code() != codes.Unauthenticated {
			t.Errorf("Expected request status to be %v and not %v", codes.Unauthenticated, reqStatus.Code())
		}
	})
}
