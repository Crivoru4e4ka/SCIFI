package handlers

import (
	"encoding/json"
	"net/http"
	"project-MVP/models"
	"project-MVP/services"
)

// GET /hypotheses
func GetHypotheses(w http.ResponseWriter, r *http.Request) {
	hypotheses, err := services.GetAllHypotheses()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hypotheses)
}

// POST /hypotheses
func CreateHypothesis(w http.ResponseWriter, r *http.Request) {
	var h models.Hypothesis
	if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	created, err := services.CreateHypothesis(h)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}
