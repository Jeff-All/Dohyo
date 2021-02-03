package authentication

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Jeff-All/Dohyo/helpers"
	"github.com/form3tech-oss/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var auth Authenticator

// Authenticator - Provides authentication
type Authenticator struct {
	log        *logrus.Logger
	config     *viper.Viper
	aud        string
	domain     string
	jwksOrigin string
}

type jwks struct {
	Keys []jsonWebKeys `json:"keys"`
}

type jsonWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// SetNewAuthenticator - Instantiates and sets a new Authenticator
func SetNewAuthenticator(
	log *logrus.Logger,
	config string,
) error {
	log.Info("setting new authenticator")
	var dir, name, ext = helpers.SplitFileName(config)

	a := Authenticator{
		log:    log,
		config: viper.New(),
	}
	a.config.AddConfigPath(dir)
	a.config.SetConfigName(name)
	a.config.SetConfigType(ext)
	if err := a.config.ReadInConfig(); err != nil {
		return err
	}

	a.aud = a.config.GetString("identifier")
	a.domain = a.config.GetString("domain")
	a.jwksOrigin = a.config.GetString("jwksOrigin")

	auth = a
	return nil
}

// AuthenticateJWT - Authenticates the provided JWT
func AuthenticateJWT(token *jwt.Token) (interface{}, error) {
	auth.log.Info("authenticating JWT")
	checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(auth.aud, false)
	if !checkAud {
		return token, errors.New("invalid audience")
	}
	checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(auth.domain, false)
	if !checkIss {
		return token, errors.New("invalid issuer")
	}

	cert, err := getPemCert(token)
	if err != nil {
		panic(err.Error())
	}

	result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	return result, nil
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(auth.jwksOrigin)

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("enable to find appropriate key")
		return cert, err
	}

	return cert, nil
}
