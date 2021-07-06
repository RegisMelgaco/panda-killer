package auth

import "time"

type Claims struct {
	Authorized bool
	AccountID  int
	Expiration time.Time
}

func NewClaims(accountID int) Claims {
	return Claims{
		Authorized: true,
		AccountID:  accountID,
		Expiration: time.Now().Add(time.Minute * 15),
	}
}
