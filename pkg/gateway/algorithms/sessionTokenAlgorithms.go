package algorithms

import (
	"local/panda-killer/cmd/config"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/auth"

	"github.com/dgrijalva/jwt-go"
)

type SessionTokenAlgorithmsImpl struct{}

func (a SessionTokenAlgorithmsImpl) GenerateSessionToken(sessionAccount *account.Account) (string, error) {
	accessSecret, err := config.GetAccessSecret()
	if err != nil {
		return "", err
	}

	claims := ToMapClaims(
		auth.NewClaims(sessionAccount.ID),
	)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedTokenString, err := token.SignedString([]byte(accessSecret))
	if err != nil {
		return "", err
	}

	return signedTokenString, nil
}

func ToMapClaims(c auth.Claims) jwt.Claims {
	return jwt.MapClaims{
		"authorized": c.Authorized,
		"account_id": c.AccountID,
		"exp":        c.Expiration.Unix(),
	}
}
