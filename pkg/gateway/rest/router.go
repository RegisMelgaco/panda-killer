package rest

import (
	"local/panda-killer/cmd/config"
	"local/panda-killer/pkg/domain/usecase"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func CreateRouter(accountUsecase *usecase.AccountUsecase, transferUsecase *usecase.TransferUsecase, authUsecase *usecase.AuthUsecase) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.SetHeader("Content-type", "application/json"))
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", CreateAccount(accountUsecase))
		r.Get("/", GetAccounts(accountUsecase))
		r.Get("/{accountID}/balance", GetAccountBalance(accountUsecase))
	})

	r.Route("/transfers", func(r chi.Router) {
		r.Use(JwtAuthentication)

		r.Post("/", CreateTransfer(transferUsecase))
		r.Get("/", ListTransfers(transferUsecase))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", Login(authUsecase))
	})

	if config.GetDebugMode() {
		restAddr, err := config.GetRestApiPort()
		if err != nil {
			panic(err)
		}
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.SetHeader("Content-type", ""))
			r.Get("/swagger/*", httpSwagger.Handler(
				httpSwagger.URL(restAddr+"/swagger/doc.json"), //The url pointing to API definition"
			))
		})
	}

	return r
}
