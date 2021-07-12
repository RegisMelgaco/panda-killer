package rest

import (
	"context"
	"fmt"
	"local/panda-killer/cmd/config"
	"local/panda-killer/pkg/domain/entity/auth"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token, err := jwt.Parse(r.Header.Get("Authorization"), func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnexpectedSigningMethod
			}
			accessSecret, err := config.GetAccessSecret()
			if err != nil {
				return []byte{}, err
			}
			return []byte(accessSecret), nil
		})

		if err != nil || !token.Valid {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		mapClaims := token.Claims.(jwt.MapClaims)

		accountID, err := strconv.ParseInt(fmt.Sprintf("%.f", mapClaims["account_id"]), 10, 64)
		if err != nil {
			panic(err)
		}
		exp, err := strconv.ParseInt(fmt.Sprintf("%.f", mapClaims["account_id"]), 10, 64)
		if err != nil {
			panic(err)
		}

		claims := auth.Claims{
			AccountID:  int(accountID),
			Authorized: mapClaims["authorized"].(bool),
			Expiration: time.Unix(exp, 0),
		}

		ctx := context.WithValue(r.Context(), auth.SessionContextKey, claims)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
