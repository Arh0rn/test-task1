package middlewares

import (
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"test-task1/pkg/logger"
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
		slog.InfoContext(
			r.Context(),
			"Request processed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
			slog.String("duration", duration.String()),
			slog.Int("status_code", rwl.StatusCode),
			slog.Int("body_size", rwl.BodySize),
			slog.String("body", string(rwl.Body)),
		)
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
