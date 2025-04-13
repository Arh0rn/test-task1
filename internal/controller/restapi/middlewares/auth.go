package middlewares

import (
	"context"
	"github.com/Arh0rn/test-task1/internal/controller/restapi/rest_errors"
	"github.com/Arh0rn/test-task1/pkg/jwtoken"
	"github.com/Arh0rn/test-task1/pkg/logger"
	"log/slog"
	"net/http"
	"strconv"
)

func AuthMiddleware(secret string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token, err := jwtoken.ExtractTokenFromRequest(r)
			if err != nil {
				rest_errors.HandleError(w, rest_errors.ErrUserUnauthorized, http.StatusUnauthorized)
				return
			}

			id, err := jwtoken.ParseToken(token, []byte(secret))
			if err != nil {
				rest_errors.HandleError(w, rest_errors.ErrUserUnauthorized, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "id", id)
			ctx = logger.WithLogUserID(ctx, strconv.Itoa(id)) //To set to every log message
			slog.InfoContext(ctx, "User authenticated")
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
