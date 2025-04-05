package middlewares

import (
	"context"
	"net/http"
	"test-task1/internal/controller/restapi/rest_errors"
	"test-task1/pkg/jwt_token"
)

func AuthMiddleware(secret string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := jwt_token.ExtractTokenFromRequest(r)
			if err != nil {
				rest_errors.HandleError(w, rest_errors.ErrUserUnauthorized, http.StatusUnauthorized)
				return
			}

			id, err := jwt_token.ParseToken(token, []byte(secret))
			if err != nil {
				rest_errors.HandleError(w, rest_errors.ErrUserUnauthorized, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "id", id)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
