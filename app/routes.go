package app

import (
	"net/http"

	"github.com/Jeff-All/Dohyo/middlewares"
	"github.com/gorilla/mux"
)

func defineRoutes(r *mux.Router) *mux.Router {
	bslog.Info("defining routes")

	r.Handle("/", routeHandlers["index"])

	r.Handle("/teams",
		middlewares.LoggingMiddleware(log,
			middlewares.CORSMiddleware(
				middleware["authorization"].BuildHandler(
					middlewares.UserMiddleware(log, db,
						routeHandlers["teams"])))))

	r.Handle("/rikishis",
		middlewares.LoggingMiddleware(log,
			middlewares.CORSMiddleware(
				middleware["authorization"].BuildHandler(
					routeHandlers["rikishis"]))))

	r.Handle("/rikishis/categorized",
		middlewares.LoggingMiddleware(log,
			middlewares.CORSMiddleware(
				middleware["authorization"].BuildHandler(
					routeHandlers["categorizedRikishis"]))))

	r.PathPrefix("/private/").Handler(
		middlewares.LoggingMiddleware(log,
			middlewares.CORSMiddleware(
				middleware["authorization"].BuildHandler(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						log.Info("handling call to the '/private' route")
						w.Write([]byte("{'message': 'SHH! It's private in here'}"))
						w.WriteHeader(http.StatusOK)
					})))))

	return r
}
