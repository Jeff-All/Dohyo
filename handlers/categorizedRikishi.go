package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Jeff-All/Dohyo/models"
	"github.com/Jeff-All/Dohyo/services"
)

// CategorizedRikishiHandler - Handles calls to the index '/'
type CategorizedRikishiHandler struct {
	Handler
	CategoryService services.CategoryService
}

func (h CategorizedRikishiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Log.Infof("CategorizeRikishiHandler serving request")

	var rikishiMap map[string][]models.Rikishi
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
	w.WriteHeader(http.StatusOK)
}

// type rikishi struct {
// 	Name   string `json:"Name"`
// 	Rank   string
// 	Avatar string
// }

// type categorizedRikishi map[string][]rikishi

// func defaultCategorizedRikishi() categorizedRikishi {
// 	return categorizedRikishi{
// 		"A": []rikishi{
// 			{Name: "Rikishi_A", Rank: "Ozeki", Avatar: "/assets/default_avatar.jpg"},
// 			{Name: "Rikishi_B", Rank: "Ozeki", Avatar: "/assets/default_avatar.jpg"},
// 			{Name: "Rikishi_C", Rank: "Ozeki", Avatar: "/assets/default_avatar.jpg"},
// 		},
// 		"B": []rikishi{
// 			{Name: "Rikishi_D", Rank: "Ozeki", Avatar: "/assets/default_avatar.jpg"},
// 			{Name: "Rikishi_E", Rank: "Ozeki", Avatar: "/assets/default_avatar.jpg"},
// 			{Name: "Rikishi_F", Rank: "Ozeki", Avatar: "/assets/default_avatar.jpg"},
// 		},
// 		"C": []rikishi{
// 			{Name: "Rikishi_G", Rank: "Ozeki", Avatar: "/assets/default_avatar.jpg"},
// 			{Name: "Rikishi_H", Rank: "Ozeki", Avatar: "/assets/default_avatar.jpg"},
// 			{Name: "Rikishi_I", Rank: "Ozeki", Avatar: "/assets/default_avatar.jpg"},
// 		},
// 		"D": []rikishi{
// 			{Name: "Rikishi_J", Rank: "Ozeki", Avatar: "/assets/default_avatar.jpg"},
// 			{Name: "Rikishi_K", Rank: "Ozeki", Avatar: "/assets/default_avatar.jpg"},
// 			{Name: "Rikishi_L", Rank: "Ozeki", Avatar: "/assets/default_avatar.jpg"},
// 		},
// 		"E": []rikishi{
// 			{Name: "Rikishi_M", Rank: "Ozeki", Avatar: "/assets/default_avatar.jpg"},
// 			{Name: "Rikishi_N", Rank: "Ozeki", Avatar: "/assets/default_avatar.jpg"},
// 			{Name: "Rikishi_O", Rank: "Ozeki", Avatar: "/assets/default_avatar.jpg"},
// 		},
// 	}
// }
