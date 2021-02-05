package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Jeff-All/Dohyo/authentication"
	"github.com/gorilla/context"
	"github.com/sirupsen/logrus"
)

type authorizationUserInfo struct {
	sub string
}

// AuthorizationMiddleware - Middleware for authorizing bearer tokens
func AuthorizationMiddleware(log *logrus.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		var resp *http.Response
		var req *http.Request
		if req, err = http.NewRequest("GET", fmt.Sprintf("%s/userinfo", authentication.Domain()), nil); err != nil {
			log.Error("error while building request to auth0 server: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", r.Header.Get("Authorization"))
		if resp, err = http.DefaultClient.Do(req); err != nil {
			log.Error("error while authorizing request: %s", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if resp.StatusCode == http.StatusUnauthorized {
			log.Info("couldn't authorize bearer token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if resp.StatusCode != http.StatusOK {
			log.Infof("unexpected status code while authorizing bearer token: %s", resp.StatusCode)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		decoder := json.NewDecoder(resp.Body)
		var json authorizationUserInfo
		if err = decoder.Decode(&json); err != nil {
			log.Errorf("error while decoding userinfo response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		context.Set(r, "user", json.sub)
		next.ServeHTTP(w, r)
	})
}
