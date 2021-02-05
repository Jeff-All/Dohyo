package middlewares

import (
	"net/http"

	"github.com/Jeff-All/Dohyo/authentication"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/sirupsen/logrus"
)

// AuthenticationMiddleware - Authentication middleware
func AuthenticationMiddleware(log *logrus.Logger, next http.Handler) http.Handler {
	m := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: authentication.AuthenticateJWT,
		SigningMethod:       jwt.SigningMethodRS256,
	})
	return m.Handler(next)
}
