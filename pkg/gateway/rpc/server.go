package rpc

import (
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/gateway/rpc/gen"
)

type Api struct {
	gen.UnimplementedPandaKillerServer

	accountUsecase *usecase.AccountUsecase
}

func NewApi(accountUsecase *usecase.AccountUsecase) *Api {
	return &Api{
		accountUsecase: accountUsecase,
	}
}
