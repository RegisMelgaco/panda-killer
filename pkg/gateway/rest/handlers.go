package rest

import (
	"encoding/json"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/usecase"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func CreateAccount(usecase *usecase.AccountUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newAccount account.Account
		err := json.NewDecoder(r.Body).Decode(&newAccount)

		if err != nil {
			log.Debugf("Failed to parse request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = usecase.CreateAccount(r.Context(), &newAccount)
		if err == account.ErrAccountCPFShouldHaveLength11 || err == account.ErrAccountNameIsObligatory {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorHolder{Message: err.Error()})
			return
		}
		if err != nil {
			log.Errorf("AccountRepo failed to create account: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newAccount)
	}
}

func GetAccounts(usecase *usecase.AccountUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := usecase.GetAccounts(r.Context())
		if err != nil {
			log.Errorf("AccountRepo failed to get accounts: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(accounts)
	}
}

func GetAccountBalance(usecase *usecase.AccountUsecase) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		accountIDStr := chi.URLParam(r, "accountID")
		accountID, err := strconv.Atoi(accountIDStr)
		log.Debug(accountID)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		balance, err := usecase.GetBalance(r.Context(), accountID)
		if err == account.ErrAccountNotFound {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
			return
		}

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(AccountBalanceResponse{Balance: balance})
	}
}

type AccountBalanceResponse struct {
	Balance float64 `json:"balance"`
}
