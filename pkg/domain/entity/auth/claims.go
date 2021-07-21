package auth

import (
	"local/panda-killer/pkg/domain/entity/account"
	"time"
)

type Claims struct {
	Authorized bool
	AccountID  account.AccountID
	Expiration time.Time
}

func NewClaims(accountID account.AccountID) Claims {
	return Claims{
		Authorized: true,
		AccountID:  accountID,
		Expiration: time.Now().Add(time.Minute * 15),
	}
}
