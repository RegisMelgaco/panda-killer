package algorithms

import (
	"fmt"
	"local/panda-killer/cmd/config"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/auth"
	"strconv"
	"time"

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

func (a SessionTokenAlgorithmsImpl) GetClaims(token string) (*auth.Claims, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		accessSecret, err := config.GetAccessSecret()
		if err != nil {
			return []byte{}, err
		}
		return []byte(accessSecret), nil
	})

	if err != nil || !t.Valid {
		return nil, auth.ErrInvalidCredentials
	}

	mapClaims := t.Claims.(jwt.MapClaims)

	accountID, err := strconv.ParseInt(fmt.Sprintf("%.f", mapClaims["account_id"]), 10, 64)
	if err != nil {
		panic(err)
	}
	exp, err := strconv.ParseInt(fmt.Sprintf("%.f", mapClaims["account_id"]), 10, 64)
	if err != nil {
		panic(err)
	}

	claims := &auth.Claims{
		AccountID:  int(accountID),
		Authorized: mapClaims["authorized"].(bool),
		Expiration: time.Unix(exp, 0),
	}

	return claims, nil
}

func ToMapClaims(c auth.Claims) jwt.Claims {
	return jwt.MapClaims{
		"authorized": c.Authorized,
		"account_id": c.AccountID,
		"exp":        c.Expiration.Unix(),
	}
}
