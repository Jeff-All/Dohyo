package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Jeff-All/Dohyo/responses"
	"github.com/Jeff-All/Dohyo/services"
)

// RikishisHandler - Handles responses for the '/rikishis' route
type RikishisHandler struct {
	Handler
	RikishiService services.RikishiService
}

// ServeHTTP - Handles calls to the route
func (h RikishisHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Log.Infof("RikishisHandler serving request")
	var err error
	var rikishis []responses.Rikishi
	if rikishis, err = h.RikishiService.GetAllCurrentCompleteRikishi(); err != nil {
		h.Log.Errorf("error while getting rikishis from RikishiService: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var arr []byte
	if arr, err = json.Marshal(rikishis); err != nil {
		h.Log.Errorf("error while parsing rikishis into json string: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(arr)
}
