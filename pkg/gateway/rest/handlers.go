package rest

import (
	"encoding/json"
	"errors"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/auth"
	"local/panda-killer/pkg/domain/entity/transfer"
	"local/panda-killer/pkg/domain/usecase"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func CreateAccount(usecase *usecase.AccountUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody CreateAccountRequest
		err := json.NewDecoder(r.Body).Decode(&requestBody)

		if err != nil {
			log.Debugf("Failed to parse request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		createdAccount, err := usecase.CreateAccount(r.Context(), requestBody.Balance, requestBody.Name, requestBody.CPF, requestBody.Password)
		if errors.Is(err, account.ErrAccountCPFShouldHaveLength11) ||
			errors.Is(err, account.ErrAccountNameIsObligatory) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
			return
		}
		if err != nil {
			log.Errorf("AccountRepo failed to create account: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CreatedAccountResponse{createdAccount.ID})
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

		var responseObjs []GetAccountResponse
		for _, a := range accounts {
			responseObjs = append(responseObjs, GetAccountResponse{ID: a.ID, Name: a.Name, CPF: a.CPF})
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseObjs)
	}
}

func GetAccountBalance(usecase *usecase.AccountUsecase) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		accountIDStr := chi.URLParam(r, "accountID")
		accountID, err := strconv.Atoi(accountIDStr)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		balance, err := usecase.GetBalance(r.Context(), accountID)
		if errors.Is(err, account.ErrAccountNotFound) {
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

func CreateTransfer(transferUsecase *usecase.TransferUsecase) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var body CreateTransferRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		createdTransfer, err := transferUsecase.CreateTransfer(
			r.Context(),
			body.OriginAccountID,
			body.DestinationAccountID,
			body.Amount,
		)
		if errors.Is(err, transfer.ErrInsufficientFundsToMakeTransaction) ||
			errors.Is(err, transfer.ErrTransferAmountShouldBeGreatterThanZero) ||
			errors.Is(err, account.ErrAccountNotFound) ||
			errors.Is(err, transfer.ErrTransferOriginAndDestinationNeedToBeDiffrent) {
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(ErrorResponse{Message: err.Error()})
			return
		}
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		json.NewEncoder(rw).Encode(CreateTransferResponse{ID: createdTransfer.ID})
	}
}

func ListTransfers(u *usecase.TransferUsecase) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		session := r.Context().Value(auth.SessionContextKey).(auth.Claims)

		transfers, err := u.ListTransfers(r.Context(), session.AccountID)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)

		var response []GetTransferResponse
		for _, t := range transfers {
			r := GetTransferResponse{
				ID:                   t.ID,
				Amount:               t.Amount,
				OriginAccountID:      t.OriginAccount.ID,
				DestinationAccountID: t.DestinationAccount.ID,
				CreatedAt:            t.CreatedAt,
			}
			response = append(response, r)
		}

		json.NewEncoder(rw).Encode(response)
	}
}

func Login(authUsecase *usecase.AuthUsecase) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var credentials LoginRequest
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		session, err := authUsecase.Login(r.Context(), credentials.CPF, credentials.Password)
		if errors.Is(err, auth.ErrInvalidCredentials) {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Authorization", session)
		rw.WriteHeader(http.StatusOK)
	}
}
