package gateway

import (
	"local/panda-killer/pkg/gateway/repository"
)

type Server struct {
	AccountRepo *repository.AccountRepo
}

func NewServer() *Server {
	return &Server{
		AccountRepo: repository.NewAccountRepo(),
	}
}
