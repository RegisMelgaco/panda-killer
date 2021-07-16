package rest

import (
	"local/panda-killer/pkg/domain/usecase"
	"net/http"
)

func JwtAuthentication(authUsecase *usecase.AuthUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx, err := authUsecase.AddClaimsToContext(r.Context(), r.Header.Get("Authorization"))
			if err != nil {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
