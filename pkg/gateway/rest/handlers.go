package rest

import (
	"encoding/json"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/usecase"
	"net/http"

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

		err = usecase.CreateAccount(&newAccount)
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
		accounts, err := usecase.GetAccounts()
		if err != nil {
			log.Errorf("AccountRepo failed to get accounts: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(accounts)
	}
}
