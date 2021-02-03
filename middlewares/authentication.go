package middlewares

import (
	"github.com/Jeff-All/Dohyo/authentication"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/urfave/negroni"
)

// BuildAuthenticationMiddleware - Builds the authentication middleware
func BuildAuthenticationMiddleware() negroni.HandlerFunc {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: authentication.AuthenticateJWT,
		SigningMethod:       jwt.SigningMethodRS256,
	}).HandlerWithNext
}
