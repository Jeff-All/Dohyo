package app

import (
	"net/http"

	"github.com/Jeff-All/Dohyo/middlewares"
	"github.com/gorilla/mux"
)

func defineRoutes(r *mux.Router) *mux.Router {
	bslog.Info("defining routes")

	r.Handle("/", routeHandlers["index"])

	r.PathPrefix("/private/").Handler(
		middlewares.LoggingMiddleware(log,
			middlewares.CORSMiddleware(
				middlewares.AuthorizationMiddleware(log,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						log.Info("handling call to the '/private' route")
						w.Write([]byte("{'message': 'SHH! It's private in here'}"))
						w.WriteHeader(http.StatusOK)
					})))))

	return r
}
