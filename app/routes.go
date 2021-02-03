package app

import (
	"net/http"

	"github.com/Jeff-All/Dohyo/middlewares"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func defineRoutes(r *mux.Router) *mux.Router {
	bslog.Info("defining routes")

	r.Handle("/", negroni.New(
		negroni.Wrap(routeHandlers["index"]),
	))

	r.PathPrefix("/private/").Handler(negroni.New(
		middlewares.BuildAuthenticationMiddleware(),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("handling call to the '/private' route")
			w.Write([]byte("SHH! It's private in here"))
		}))))

	return r
}
