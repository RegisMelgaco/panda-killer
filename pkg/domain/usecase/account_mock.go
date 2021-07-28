// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package usecase

import (
	"context"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/shared"
	"sync"
)

// Ensure, that AccountUsecaseMock does implement AccountUsecase.
// If this is not the case, regenerate this file with moq.
var _ AccountUsecase = &AccountUsecaseMock{}

// AccountUsecaseMock is a mock implementation of AccountUsecase.
//
// 	func TestSomethingThatUsesAccountUsecase(t *testing.T) {
//
// 		// make and configure a mocked AccountUsecase
// 		mockedAccountUsecase := &AccountUsecaseMock{
// 			CreateAccountFunc: func(ctx context.Context, balance shared.Money, name string, cpf string, password string) (*account.Account, error) {
// 				panic("mock out the CreateAccount method")
// 			},
// 			GetAccountsFunc: func(contextMoqParam context.Context) ([]account.Account, error) {
// 				panic("mock out the GetAccounts method")
// 			},
// 			GetBalanceFunc: func(contextMoqParam context.Context, accountID account.AccountID) (shared.Money, error) {
// 				panic("mock out the GetBalance method")
// 			},
// 		}
//
// 		// use mockedAccountUsecase in code that requires AccountUsecase
// 		// and then make assertions.
//
// 	}
type AccountUsecaseMock struct {
	// CreateAccountFunc mocks the CreateAccount method.
	CreateAccountFunc func(ctx context.Context, balance shared.Money, name string, cpf string, password string) (*account.Account, error)

	// GetAccountsFunc mocks the GetAccounts method.
	GetAccountsFunc func(contextMoqParam context.Context) ([]account.Account, error)

	// GetBalanceFunc mocks the GetBalance method.
	GetBalanceFunc func(contextMoqParam context.Context, accountID account.AccountID) (shared.Money, error)

	// calls tracks calls to the methods.
	calls struct {
		// CreateAccount holds details about calls to the CreateAccount method.
		CreateAccount []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Balance is the balance argument value.
			Balance shared.Money
			// Name is the name argument value.
			Name string
			// Cpf is the cpf argument value.
			Cpf string
			// Password is the password argument value.
			Password string
		}
		// GetAccounts holds details about calls to the GetAccounts method.
		GetAccounts []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
		}
		// GetBalance holds details about calls to the GetBalance method.
		GetBalance []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// AccountID is the accountID argument value.
			AccountID account.AccountID
		}
	}
	lockCreateAccount sync.RWMutex
	lockGetAccounts   sync.RWMutex
	lockGetBalance    sync.RWMutex
}

// CreateAccount calls CreateAccountFunc.
func (mock *AccountUsecaseMock) CreateAccount(ctx context.Context, balance shared.Money, name string, cpf string, password string) (*account.Account, error) {
	if mock.CreateAccountFunc == nil {
		panic("AccountUsecaseMock.CreateAccountFunc: method is nil but AccountUsecase.CreateAccount was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		Balance  shared.Money
		Name     string
		Cpf      string
		Password string
	}{
		Ctx:      ctx,
		Balance:  balance,
		Name:     name,
		Cpf:      cpf,
		Password: password,
	}
	mock.lockCreateAccount.Lock()
	mock.calls.CreateAccount = append(mock.calls.CreateAccount, callInfo)
	mock.lockCreateAccount.Unlock()
	return mock.CreateAccountFunc(ctx, balance, name, cpf, password)
}

// CreateAccountCalls gets all the calls that were made to CreateAccount.
// Check the length with:
//     len(mockedAccountUsecase.CreateAccountCalls())
func (mock *AccountUsecaseMock) CreateAccountCalls() []struct {
	Ctx      context.Context
	Balance  shared.Money
	Name     string
	Cpf      string
	Password string
} {
	var calls []struct {
		Ctx      context.Context
		Balance  shared.Money
		Name     string
		Cpf      string
		Password string
	}
	mock.lockCreateAccount.RLock()
	calls = mock.calls.CreateAccount
	mock.lockCreateAccount.RUnlock()
	return calls
}

// GetAccounts calls GetAccountsFunc.
func (mock *AccountUsecaseMock) GetAccounts(contextMoqParam context.Context) ([]account.Account, error) {
	if mock.GetAccountsFunc == nil {
		panic("AccountUsecaseMock.GetAccountsFunc: method is nil but AccountUsecase.GetAccounts was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
	}{
		ContextMoqParam: contextMoqParam,
	}
	mock.lockGetAccounts.Lock()
	mock.calls.GetAccounts = append(mock.calls.GetAccounts, callInfo)
	mock.lockGetAccounts.Unlock()
	return mock.GetAccountsFunc(contextMoqParam)
}

// GetAccountsCalls gets all the calls that were made to GetAccounts.
// Check the length with:
//     len(mockedAccountUsecase.GetAccountsCalls())
func (mock *AccountUsecaseMock) GetAccountsCalls() []struct {
	ContextMoqParam context.Context
} {
	var calls []struct {
		ContextMoqParam context.Context
	}
	mock.lockGetAccounts.RLock()
	calls = mock.calls.GetAccounts
	mock.lockGetAccounts.RUnlock()
	return calls
}

// GetBalance calls GetBalanceFunc.
func (mock *AccountUsecaseMock) GetBalance(contextMoqParam context.Context, accountID account.AccountID) (shared.Money, error) {
	if mock.GetBalanceFunc == nil {
		panic("AccountUsecaseMock.GetBalanceFunc: method is nil but AccountUsecase.GetBalance was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		AccountID       account.AccountID
	}{
		ContextMoqParam: contextMoqParam,
		AccountID:       accountID,
	}
	mock.lockGetBalance.Lock()
	mock.calls.GetBalance = append(mock.calls.GetBalance, callInfo)
	mock.lockGetBalance.Unlock()
	return mock.GetBalanceFunc(contextMoqParam, accountID)
}

// GetBalanceCalls gets all the calls that were made to GetBalance.
// Check the length with:
//     len(mockedAccountUsecase.GetBalanceCalls())
func (mock *AccountUsecaseMock) GetBalanceCalls() []struct {
	ContextMoqParam context.Context
	AccountID       account.AccountID
} {
	var calls []struct {
		ContextMoqParam context.Context
		AccountID       account.AccountID
	}
	mock.lockGetBalance.RLock()
	calls = mock.calls.GetBalance
	mock.lockGetBalance.RUnlock()
	return calls
}