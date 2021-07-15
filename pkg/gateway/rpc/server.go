package rpc

import (
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/gateway/rpc/gen"

	"google.golang.org/grpc"
)

type Api struct {
	gen.UnimplementedPandaKillerServer

	AccountUsecase *usecase.AccountUsecase
	AuthUsecase    *usecase.AuthUsecase
}

func NewApi(accountUsecase *usecase.AccountUsecase, authUsecase *usecase.AuthUsecase) *Api {
	return &Api{
		AccountUsecase: accountUsecase,
		AuthUsecase:    authUsecase,
	}
}

func NewServer(api *Api) *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(SessionInterceptor(*api.AuthUsecase)),
	)
	gen.RegisterPandaKillerServer(s, api)
	return s
}
