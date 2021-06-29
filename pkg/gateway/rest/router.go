package rest

import (
	"local/panda-killer/pkg/domain/usecase"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func CreateRouter(accountUsecase *usecase.AccountUsecase) http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/accounts/", CreateAccount(accountUsecase))
	r.Get("/accounts/", GetAccounts(accountUsecase))
	r.Get("/accounts/{accountID}/balance", GetAccountBalance(accountUsecase))

	return r
}
