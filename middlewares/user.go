package middlewares

import (
	"net/http"

	"github.com/Jeff-All/Dohyo/models"
	"github.com/gorilla/context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// UserMiddleware - Pulls the user model by the Auth0 user ID
func UserMiddleware(log *logrus.Logger, db *gorm.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID interface{}
		var userIDString string
		var ok bool
		if userID = context.Get(r, "userID"); userID == nil {
			log.Errorf("userID has not been defined in this request's context")
			return
		} else if userIDString, ok = userID.(string); !ok {
			log.Errorf("userID defined in this context must be of type string")
			return
		}
		var email interface{}
		var emailString string
		if email = context.Get(r, "email"); email == nil {
			log.Errorf("user email has not been defined in this request's context")
			return
		} else if emailString, ok = email.(string); !ok {
			log.Errorf("email defined in this context must be of type string")
			return
		}

		user := models.User{}
		db.Where(models.User{Auth0ID: userIDString}).Attrs(models.User{Email: emailString}).FirstOrCreate(&user)

		context.Set(r, "user", user)

		next.ServeHTTP(w, r)
	})
}
