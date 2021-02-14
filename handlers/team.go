package handlers

import "net/http"

// TeamHandler - Handles calls to the '/team' routes
type TeamHandler struct {
	Handler
}

func (h TeamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Log.Infof("TeamHandler serving request")

}
