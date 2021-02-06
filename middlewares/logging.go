package middlewares

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// LoggingMiddleware - Logs the route
func LoggingMiddleware(log *logrus.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Handling '%s' request for '%s'", r.Method, r.RequestURI)

		next.ServeHTTP(w, r)
	})
}
