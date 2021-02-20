package middlewares

import (
	"net/http"
)

// MiddlewareInterface - Standard Middleware interface for building handlers
type MiddlewareInterface interface {
	BuildHandler(next http.Handler) http.Handler
}
