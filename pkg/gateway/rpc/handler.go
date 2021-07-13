package rpc

import (
	"context"
	"errors"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/gateway/rpc/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Api) CreateAccount(ctx context.Context, accountReq *gen.CreateAccountRequest) (*gen.CreateAccountResponse, error) {
	createdAccount, err := s.accountUsecase.CreateAccount(
		ctx, int(accountReq.Balance), accountReq.Name, accountReq.Cpf, accountReq.Password,
	)

	if errors.Is(err, account.ErrAccountCPFShouldHaveLength11) ||
		errors.Is(err, account.ErrAccountNameIsObligatory) ||
		errors.Is(err, account.ErrAccountCPFShouldBeUnique) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}

	return &gen.CreateAccountResponse{Id: int32(createdAccount.ID)}, nil
}

func (s *Api) GetAccountBalance(ctx context.Context, request *gen.GetAccountBalanceRequest) (*gen.GetAccountBalanceResponse, error) {
	balance, err := s.accountUsecase.GetBalance(ctx, int(request.AccountId))
	if errors.Is(err, account.ErrAccountNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}

	return &gen.GetAccountBalanceResponse{
		Balance: int32(balance),
	}, nil
}
