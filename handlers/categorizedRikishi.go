package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Jeff-All/Dohyo/services"
)

// CategorizedRikishiHandler - Handles calls to the index '/'
type CategorizedRikishiHandler struct {
	Handler
	CategoryService services.CategoryService
}

func (h CategorizedRikishiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Log.Infof("CategorizeRikishiHandler serving request")

	var rikishiMap map[string][]string
	var err error
	if rikishiMap, err = h.CategoryService.GetRikishiByCategory(); err != nil {
		h.Log.Errorf("error while retriving rikishi by category: %s", err)
	}
	h.Log.Infof("map: %s", rikishiMap)

	if err = json.NewEncoder(w).Encode(rikishiMap); err != nil {
		h.Log.Errorf("error encoding response for CategorizedRikishiService: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}
