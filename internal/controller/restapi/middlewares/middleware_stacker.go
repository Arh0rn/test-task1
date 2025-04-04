package middlewares

import "net/http"

type Middleware func(next http.Handler) http.Handler

func CreateMiddlewareStack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, middleware := range middlewares {
			next = middleware(next)
		}
		return next
	}
}
