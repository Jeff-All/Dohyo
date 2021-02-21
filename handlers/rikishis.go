package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Jeff-All/Dohyo/models"
	"github.com/Jeff-All/Dohyo/services"
)

type rikishiResponse struct {
	ID      uint
	Name    string
	Avatar  string
	Rank    string
	Wins    uint
	Loss    uint
	Results []resultResponse
	Matches []matchResponse
}

type resultResponse struct {
	Tournament string
	Wins       uint
	Losses     uint
}

type matchResponse struct {
	Day      uint
	Opponent uint
	Result   string
}

// RikishisHandler - Handles responses for the '/rikishis' route
type RikishisHandler struct {
	Handler
	RikishiService services.RikishiService
}

// ServeHTTP - Handles calls to the route
func (h RikishisHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Log.Infof("RikishisHandler serving request")
	var err error
	var rikishis []models.Rikishi
	if rikishis, err = h.RikishiService.GetAllRikishi(); err != nil {
		h.Log.Errorf("error while getting rikishis from RikishiService: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respRikishis := make([]rikishiResponse, len(rikishis))
	for index, cur := range rikishis {
		respRikishis[index] = rikishiResponse{
			ID:     cur.ID,
			Name:   cur.Name,
			Avatar: cur.Avatar,
			Rank:   cur.Rank,
			Wins:   11,
			Loss:   4,
		}
	}

	var arr []byte
	if arr, err = json.Marshal(respRikishis); err != nil {
		h.Log.Errorf("error while parsing rikishis into json string: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(arr)
}
