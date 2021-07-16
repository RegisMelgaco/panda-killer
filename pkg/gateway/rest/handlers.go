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

// CreateAccount handles account creation requests
// @Summary Create account
// @Description creates a new account with provided info
// @Tags Account
//
// @Accept  json
// @Produce  json
// @Param account body rest.CreateAccountRequest true "info used in account creation"
//
// @Success 201 {object} CreatedAccountResponse
// @Failure 400 {object} ErrorResponse "Possible errors: account.ErrAccountCPFShouldHaveLength11, account.ErrAccountNameIsObligatory and account.ErrAccountCPFShouldBeUnique"
// @Failure 500
// @Router /accounts [post]
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
			errors.Is(err, account.ErrAccountNameIsObligatory) ||
			errors.Is(err, account.ErrAccountCPFShouldBeUnique) {
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

// GetAccounts handles account listing requests
// @Summary List accounts
// @Description Lists all created accounts
// @Tags Account
//
// @Produce  json
//
// @Success 200 {object} []GetAccountResponse
// @Failure 500
// @Router /accounts [get]
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

// GetAccountBalance handles account balance consultation requests
// @Summary Get account balance
// @Description Get the account's balance
// @Tags Account
//
// @Produce  json
// @Param accountID path int true "ID from the account owner of the desired balance"
//
// @Success 200 {object} AccountBalanceResponse
// @Failure 400 {string} string "accountID not in int format"
// @Failure 404 {string} string "Account not found"
// @Failure 500
// @Router /{accountID}/balance [get]
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

// CreateTransfer handles transfer creation requests
// @Summary Transfer account
// @Description creates a new transfer from origin account to destination account with desired amount
// @Tags Transfer
//
// @Accept  json
// @Produce  json
// @Param transfer body rest.CreateTransferRequest true "Contains the origin account (source of the money), destination account and amount to be transferred."
//
// @Success 201 {object} CreateTransferResponse
// @Failure 400 {object} ErrorResponse "It contains the error reason. Possible errors (ErrInsufficientFundsToMakeTransaction, ErrTransferAmountShouldBeGreatterThanZero, ErrAccountNotFoundErrTransferOriginAndDestinationNeedToBeDiffrent)"
// @Failure 500
// @Router /transfers [post]
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

// ListTransfers handles listing transfers from a account requests
// @Summary List account transfers
// @Description Lists all transfers where the logged account takes part
// @Tags Transfer
//
// @Produce  json
// @Security ApiKeyAuth
//
// @Success 200 {object} []GetAccountResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 500
// @Router /transfers [get]
func ListTransfers(u *usecase.TransferUsecase) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		session := r.Context().Value(auth.SessionContextKey).(*auth.Claims)

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

// Login handles login requests
// @Summary Login
// @Description Login a user based on cpf and password from a created Account
// @Tags Auth
//
// @Accept  json
// @Param transfer body rest.LoginRequest true "Login credentials"
//
// @Success 200 {object} CreateTransferResponse
// @Failure 401 {string} string "It was not possible to find a account with the informed cpf or the password doesn't match to the secret."
// @Failure 500
// @Router /auth/login [post]
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
