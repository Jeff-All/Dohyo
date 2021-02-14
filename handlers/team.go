package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Jeff-All/Dohyo/models"
	"github.com/Jeff-All/Dohyo/services"
	"github.com/gorilla/context"
)

// TeamHandler - Handles calls to the '/team' routes
type TeamHandler struct {
	Handler
	CategoryService services.CategoryService
	TeamService     services.TeamService
}

func (h TeamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Log.Infof("TeamHandler serving request")

	// Retrieve User from context
	userI := context.Get(r, "user")
	if userI == nil {
		h.Log.Errorf("user is not set in the request's context")
		return
	}
	var user models.User
	var ok bool
	if user, ok = userI.(models.User); !ok {
		h.Log.Errorf("user needs to be a models.User")
		return
	}

	switch r.Method {
	case "PUT":
		h.put(user, w, r)
		break
	default:
		h.Log.Infof("invalid method %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (h TeamHandler) put(user models.User, w http.ResponseWriter, r *http.Request) {
	h.Log.Infof("serving PUT")
	decoder := json.NewDecoder(r.Body)
	requestBody := map[string]uint{}
	var err error
	if err = decoder.Decode(&requestBody); err != nil {
		h.Log.Errorf("error while decoding PUT request: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var categoryCount int
	if categoryCount, err = h.CategoryService.Count(); err != nil {
		h.Log.Errorf("error pulling category count: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if len(requestBody) != categoryCount {
		h.Log.Infof("request body needs %d entries but has %d", categoryCount, len(requestBody))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	index := 0
	rikishis := make([]models.Rikishi, categoryCount)
	for _, rikishiID := range requestBody {
		h.Log.Infof("rikishis: %d", rikishiID)
		rikishis[index].ID = rikishiID
		index++
	}
	var distinctCategoryCount int
	if distinctCategoryCount, err = h.CategoryService.GetCategoryCountOfRikishis(rikishis); err != nil {
		h.Log.Errorf("Error while verfiying rikishi categories: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if distinctCategoryCount != categoryCount {
		h.Log.Infof("rikishi must all be in differenct categories")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.TeamService.SaveRikishisToTeam(user, rikishis); err != nil {
		h.Log.Errorf("error while saving rikishi to team: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	return
}
