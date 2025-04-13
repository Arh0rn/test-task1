package middlewares

import (
	"github.com/Arh0rn/test-task1/pkg/logger"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := uuid.New()
		ctx = logger.WithLogRequestID(ctx, requestID.String())
		r = r.WithContext(ctx)

		start := time.Now()
		rwl := NewResponseLogger(w)

		next.ServeHTTP(rwl, r)

		duration := time.Since(start)

		if strings.HasPrefix(r.URL.Path, "/swagger/") {
			slog.InfoContext(
				r.Context(), "Request processed",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"duration", duration.String(),
				"status_code", rwl.StatusCode,
				"body_size", rwl.BodySize,
				//"body", string(rwl.Body), // With no body to swagger html
			)

		} else {
			slog.InfoContext(
				r.Context(), "Request processed",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"duration", duration.String(),
				"status_code", rwl.StatusCode,
				"body_size", rwl.BodySize,
				"body", string(rwl.Body),
			)
		}
	})
}

type ResponseLogger struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
	BodySize   int
}

func NewResponseLogger(w http.ResponseWriter) *ResponseLogger {
	return &ResponseLogger{ResponseWriter: w, StatusCode: http.StatusOK}
}

func (r *ResponseLogger) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseLogger) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.BodySize += size
	r.Body = append(r.Body, b...)
	return size, err
}
