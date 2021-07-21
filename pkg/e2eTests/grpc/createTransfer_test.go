package e2etest

import (
	"context"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/shared"
	"local/panda-killer/pkg/domain/entity/transfer"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/gateway/algorithms"
	"local/panda-killer/pkg/gateway/db/postgres"
	"local/panda-killer/pkg/gateway/repository"
	"local/panda-killer/pkg/gateway/rpc"
	"local/panda-killer/pkg/gateway/rpc/gen"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestCreateTransfer(t *testing.T) {
	ctx := context.Background()

	postgres.RunMigrations()
	defer postgres.DownToMigrationZero()

	pgxConn, _ := postgres.OpenConnection()
	passAlgo := algorithms.PasswordHashingAlgorithmsImpl{}
	sessionAlgo := algorithms.SessionTokenAlgorithmsImpl{}
	accountRepo := repository.NewAccountRepo(pgxConn)
	transferRepo := repository.NewTransferRepo(pgxConn)
	accountUsecase := usecase.NewAccountUsecase(accountRepo, passAlgo)
	authUsecase := usecase.NewAuthUsecase(accountRepo, sessionAlgo, passAlgo)
	transferUsecase := usecase.NewTransferUsecase(transferRepo, accountRepo)

	api := &rpc.Api{
		AuthUsecase:     authUsecase,
		TransferUsecase: transferUsecase,
	}

	client, conn, err := StartServerAndGetClient(ctx, api)
	defer conn.Close()

	passord := ";)"
	testAccount1, err := accountUsecase.CreateAccount(context.Background(), 1, "João", "85385023495", passord)
	if err != nil {
		t.Errorf("Failed to create test user: %v", err)
	}

	resp, err := client.Login(ctx, &gen.LoginRequest{
		Cpf:      testAccount1.CPF,
		Password: passord,
	})

	_, ok := status.FromError(err)
	if !ok {
		t.Errorf("Failed to request login: %v", err)
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", resp.Token)

	t.Run("Create transfer without authentication should return unauthenticated", func(t *testing.T) {
		testAccount2 := account.Account{Balance: 2, Name: "Joana", CPF: "35081960896", Secret: "s"}
		err = accountRepo.CreateAccount(context.Background(), &testAccount2)
		if err != nil {
			t.Errorf("Failed to create test account2: %v", err)
		}

		_, err := client.CreateTransfer(context.Background(), &gen.CreateTransferRequest{
			OriginAccountId:      int32(testAccount1.ID),
			DestinationAccountId: int32(testAccount2.ID),
			Amount:               1,
		})

		reqStatus, _ := status.FromError(err)
		if reqStatus.Code() != codes.Unauthenticated {
			t.Errorf("Expected response status was %v and not %v", codes.Unauthenticated, reqStatus.Code())
		}
	})

	t.Run("Create transfer with success should update users balances with success", func(t *testing.T) {
		testAccount2 := account.Account{Balance: 2, Name: "Joana", CPF: "46834635203", Secret: "s"}
		err = accountRepo.CreateAccount(context.Background(), &testAccount2)
		if err != nil {
			t.Errorf("Failed to create test account2: %v", err)
		}

		resp, err := client.CreateTransfer(ctx, &gen.CreateTransferRequest{
			OriginAccountId:      int32(testAccount1.ID),
			DestinationAccountId: int32(testAccount2.ID),
			Amount:               1,
		})

		reqStatus, _ := status.FromError(err)
		if reqStatus.Code() != codes.OK {
			t.Errorf("Transfer creation request response should be %v not %v", codes.OK, reqStatus.Code())
		}

		if resp.Id < 1 {
			t.Errorf("Response body with invalid id: %v", resp.Id)
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
		var originalOriginAccountBalance shared.Money = 1
		var originalDestineAccountBalance shared.Money = 0

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

		transferRequest := &gen.CreateTransferRequest{OriginAccountId: int32(testAccount1.ID), DestinationAccountId: int32(testAccount2.ID), Amount: int32(originalOriginAccountBalance + 1)}
		_, err = client.CreateTransfer(ctx, transferRequest)

		reqStatus, _ := status.FromError(err)
		if reqStatus.Code() != codes.InvalidArgument {
			t.Errorf("Transfer creation request response should be %v not %v", codes.InvalidArgument, reqStatus.Code())
		}

		if reqStatus.Message() != transfer.ErrInsufficientFundsToMakeTransaction.Error() {
			t.Errorf(
				"Response message should be %v and not %v",
				transfer.ErrInsufficientFundsToMakeTransaction.Error(),
				reqStatus.Message(),
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
		transferRequest := &gen.CreateTransferRequest{OriginAccountId: 132, DestinationAccountId: 13212, Amount: 1}
		_, err := client.CreateTransfer(ctx, transferRequest)

		reqStatus, _ := status.FromError(err)
		if reqStatus.Code() != codes.InvalidArgument {
			t.Errorf("Transfer creation request response should be %v not %v", codes.InvalidArgument, reqStatus.Code())
		}

		if reqStatus.Message() != account.ErrAccountNotFound.Error() {
			t.Errorf(
				"Response message should be %v and not %v",
				account.ErrAccountNotFound.Error(),
				reqStatus.Message(),
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

		_, err = client.CreateTransfer(ctx, &gen.CreateTransferRequest{
			OriginAccountId:      int32(testAccount1.ID),
			DestinationAccountId: int32(testAccount2.ID),
			Amount:               0,
		})

		reqStatus, _ := status.FromError(err)
		if reqStatus.Code() != codes.InvalidArgument {
			t.Errorf("Response status should be %v and not %v", codes.InvalidArgument, reqStatus.Code())
		}

		if reqStatus.Message() != transfer.ErrTransferAmountShouldBeGreatterThanZero.Error() {
			t.Errorf("Expected message in response was '%v' and not '%v'", transfer.ErrTransferAmountShouldBeGreatterThanZero.Error(), reqStatus.Message())
		}
	})
	t.Run("Should not be possible to transfer to your self", func(t *testing.T) {
		testAccount := account.Account{Name: "Maria", CPF: "06268075730", Secret: "s"}
		err := accountRepo.CreateAccount(context.Background(), &testAccount)
		if err != nil {
			t.Errorf("Failed to create test account1: %v", err)
		}

		_, err = client.CreateTransfer(ctx, &gen.CreateTransferRequest{
			OriginAccountId:      int32(testAccount.ID),
			DestinationAccountId: int32(testAccount.ID),
			Amount:               42,
		})

		reqStatus, _ := status.FromError(err)
		if reqStatus.Code() != codes.InvalidArgument {
			t.Errorf("Response status should be %v and not %v", codes.InvalidArgument, reqStatus.Code())
		}

		if reqStatus.Message() != transfer.ErrTransferOriginAndDestinationNeedToBeDiffrent.Error() {
			t.Errorf("Expected message in response was '%v' and not '%v'", transfer.ErrTransferOriginAndDestinationNeedToBeDiffrent.Error(), reqStatus.Message())
		}
	})
}
