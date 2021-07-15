package rpc

import (
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/gateway/rpc/gen"

	"google.golang.org/grpc"
)

type Api struct {
	gen.UnimplementedPandaKillerServer

	AccountUsecase  *usecase.AccountUsecase
	AuthUsecase     *usecase.AuthUsecase
	TransferUsecase *usecase.TransferUsecase
}

func NewApi(accountUsecase *usecase.AccountUsecase, authUsecase *usecase.AuthUsecase, transferUsecase *usecase.TransferUsecase) *Api {
	return &Api{
		AccountUsecase:  accountUsecase,
		AuthUsecase:     authUsecase,
		TransferUsecase: transferUsecase,
	}
}

func NewServer(api *Api) *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(SessionInterceptor(*api.AuthUsecase)),
	)
	gen.RegisterPandaKillerServer(s, api)
	return s
}
