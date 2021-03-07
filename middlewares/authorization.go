package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Jeff-All/Dohyo/authentication"
	"github.com/gorilla/context"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type authorizationUserInfo struct {
	Sub   string
	Email string
}

// AuthorizationMiddleware - Middleware for authorizing bearer tokens
type AuthorizationMiddleware struct {
	Log   *logrus.Logger
	Cache *cache.Cache
}

// BuildHandler - Middleware for authorizing bearer tokens
func (m *AuthorizationMiddleware) BuildHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		m.Log.Infof("authorizing '%s'", bearerToken)
		var user authorizationUserInfo
		if userI, found := m.Cache.Get(bearerToken); !found {
			m.Log.Infof("cache miss for %s", bearerToken)
			var err error
			var resp *http.Response
			var req *http.Request

			if req, err = http.NewRequest("GET", fmt.Sprintf("%s/userinfo", authentication.Domain()), nil); err != nil {
				m.Log.Error("error while building request to auth0 server: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			req.Header.Set("Authorization", r.Header.Get("Authorization"))
			if resp, err = http.DefaultClient.Do(req); err != nil {
				m.Log.Error("error while authorizing request: %s", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusUnauthorized {
				m.Log.Info("couldn't authorize bearer token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else if resp.StatusCode != http.StatusOK {
				m.Log.Infof("unexpected status code while authorizing bearer token: %s", resp.StatusCode)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			decoder := json.NewDecoder(resp.Body)
			user := authorizationUserInfo{}
			if err = decoder.Decode(&user); err != nil {
				m.Log.Errorf("error while decoding userinfo response: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			m.Cache.Set(bearerToken, &user, cache.DefaultExpiration)
		} else {
			m.Log.Infof("cache hit for %s", bearerToken)
			user = *(userI.(*authorizationUserInfo))
		}

		m.Log.Infof("userID='%s', email='%s'", user.Sub, user.Email)
		context.Set(r, "userID", user.Sub)
		context.Set(r, "email", user.Email)
		next.ServeHTTP(w, r)
	})
}
