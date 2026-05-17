package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"project-MVP/models"
	"project-MVP/services"
)

// ------------------- Grants -------------------

// POST /grants
func CreateGrant(w http.ResponseWriter, r *http.Request) {
	userID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	var g models.Grant
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	g.CreatedBy = userID

	created, err := services.CreateGrant(g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

// GET /grants
func ListGrants(w http.ResponseWriter, r *http.Request) {
	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	query := r.URL.Query()
	filters := services.ListGrantsFilters{
		Status:       strings.TrimSpace(query.Get("status")),
		GrantType:    strings.TrimSpace(query.Get("type")),
		Organization: strings.TrimSpace(query.Get("organization")),
		Search:       strings.TrimSpace(query.Get("search")),
	}

	if limitStr := query.Get("limit"); limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
			filters.Limit = v
		}
	}
	if offsetStr := query.Get("offset"); offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil && v >= 0 {
			filters.Offset = v
		}
	}
	if fromStr := query.Get("date_from"); fromStr != "" {
		if t, err := time.Parse("2006-01-02", fromStr); err == nil {
			filters.DateFrom = &t
		}
	}
	if toStr := query.Get("date_to"); toStr != "" {
		if t, err := time.Parse("2006-01-02", toStr); err == nil {
			filters.DateTo = &t
		}
	}

	grants, total, err := services.ListGrants(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"grants": grants,
		"total":  total,
	})
}

// GET /grants/{id}
func GetGrant(w http.ResponseWriter, r *http.Request) {
	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		http.Error(w, "invalid grant id", http.StatusBadRequest)
		return
	}

	grant, err := services.GetGrantByID(id)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "grant not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grant)
}

// PATCH /grants/{id}
func UpdateGrant(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		http.Error(w, "invalid grant id", http.StatusBadRequest)
		return
	}

	grant, err := services.GetGrantByID(id)
	if err != nil {
		http.Error(w, "grant not found", http.StatusNotFound)
		return
	}
	if grant.CreatedBy != currentUserID && !services.IsAdmin(currentUserID) {
		http.Error(w, "insufficient permissions", http.StatusForbidden)
		return
	}

	var g models.Grant
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := services.UpdateGrant(id, g)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "grant not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DELETE /grants/{id}
func DeleteGrant(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		http.Error(w, "invalid grant id", http.StatusBadRequest)
		return
	}

	grant, err := services.GetGrantByID(id)
	if err != nil {
		http.Error(w, "grant not found", http.StatusNotFound)
		return
	}
	if grant.CreatedBy != currentUserID && !services.IsAdmin(currentUserID) {
		http.Error(w, "insufficient permissions", http.StatusForbidden)
		return
	}

	if err := services.DeleteGrant(id); err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "grant not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /grants/{id}/projects
func GetGrantProjects(w http.ResponseWriter, r *http.Request) {
	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		http.Error(w, "invalid grant id", http.StatusBadRequest)
		return
	}

	projects, err := services.GetGrantProjects(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// GET /grants/{id}/budget
func GetGrantBudget(w http.ResponseWriter, r *http.Request) {
	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		http.Error(w, "invalid grant id", http.StatusBadRequest)
		return
	}

	grant, err := services.GetGrantByID(id)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "grant not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allocated, err := services.GetGrantTotalAllocated(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":     grant.TotalAmount,
		"allocated": allocated,
		"remaining": grant.TotalAmount - allocated,
		"currency":  grant.Currency,
	})
}

// ------------------- Project Grant Funding -------------------

// POST /grants/{id}/projects
func AddProjectToGrant(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	grantID, err := strconv.Atoi(vars["id"])
	if err != nil || grantID <= 0 {
		http.Error(w, "invalid grant id", http.StatusBadRequest)
		return
	}

	grant, err := services.GetGrantByID(grantID)
	if err != nil {
		http.Error(w, "grant not found", http.StatusNotFound)
		return
	}
	if grant.CreatedBy != currentUserID && !services.IsAdmin(currentUserID) {
		http.Error(w, "insufficient permissions", http.StatusForbidden)
		return
	}

	var f models.ProjectGrantFunding
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	f.GrantId = grantID

	created, err := services.CreateProjectGrantFunding(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

// DELETE /grants/{grantId}/projects/{fundingId}
func RemoveProjectFromGrant(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	grantID, err := strconv.Atoi(vars["grantId"])
	if err != nil || grantID <= 0 {
		http.Error(w, "invalid grant id", http.StatusBadRequest)
		return
	}

	grant, err := services.GetGrantByID(grantID)
	if err != nil {
		http.Error(w, "grant not found", http.StatusNotFound)
		return
	}
	if grant.CreatedBy != currentUserID && !services.IsAdmin(currentUserID) {
		http.Error(w, "insufficient permissions", http.StatusForbidden)
		return
	}

	fundingID, err := strconv.Atoi(vars["fundingId"])
	if err != nil || fundingID <= 0 {
		http.Error(w, "invalid funding id", http.StatusBadRequest)
		return
	}

	if err := services.DeleteProjectGrantFunding(fundingID); err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "funding not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /projects/{id}/grants
func GetProjectGrants(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil || projectID <= 0 {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	if !RequirePermission(w, r, projectID, "project.view") {
		return
	}

	grants, err := services.GetProjectGrants(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grants)
}
