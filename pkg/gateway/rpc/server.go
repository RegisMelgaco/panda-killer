package rpc

import (
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/gateway/rpc/gen"
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
