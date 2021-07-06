package algorithms

import (
	"local/panda-killer/cmd/config"
	"local/panda-killer/pkg/domain/entity/account"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type SessionTokenAlgorithmsImpl struct{}

func (a SessionTokenAlgorithmsImpl) GenerateSessionToken(sessionAccount *account.Account) (string, error) {
	accessSecret, err := config.GetAccessSecret()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"authorized": true,
		"account_id": sessionAccount.ID,
		"exp":        time.Now().Add(time.Minute * 15).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedTokenString, err := token.SignedString([]byte(accessSecret))
	if err != nil {
		return "", err
	}

	return signedTokenString, nil
}
