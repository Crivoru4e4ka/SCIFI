package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"project-MVP/models"
	"project-MVP/services"

	"github.com/gorilla/mux"
)

// ------------------- Grants -------------------

// CreateGrant godoc
// @Summary Создать новый грант
// @Description Создает запись о гранте или источнике финансирования
// @Tags grants
// @Accept json
// @Produce json
// @Param grant body models.Grant true "Данные гранта"
// @Success 200 {object} models.Grant
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /grants [post]
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

// ListGrants godoc
// @Summary Список грантов
// @Description Получить список всех грантов с возможностью фильтрации по статусу, типу и поисковому запросу
// @Tags grants
// @Param status query string false "Фильтр по статусу (draft, active, completed...)"
// @Param type query string false "Фильтр по типу (state, international...)"
// @Param search query string false "Поиск по названию или организации"
// @Param limit query int false "Лимит записей"
// @Param offset query int false "Смещение"
// @Produce json
// @Success 200 {object} map[string]interface{} "Объект с массивом 'grants' и числом 'total'"
// @Failure 500 {string} string "Internal Server Error"
// @Router /grants [get]
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

// GetGrant godoc
// @Summary Детальная информация о гранте
// @Tags grants
// @Produce json
// @Param id path int true "ID гранта"
// @Success 200 {object} models.Grant
// @Failure 404 {string} string "grant not found"
// @Router /grants/{id} [get]
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

// UpdateGrant godoc
// @Summary Обновить данные гранта
// @Description Изменить описание, сроки или сумму гранта. Доступно автору или админу.
// @Tags grants
// @Accept json
// @Produce json
// @Param id path int true "ID гранта"
// @Param grant body models.Grant true "Новые данные"
// @Success 200 {object} models.Grant
// @Failure 403 {string} string "insufficient permissions"
// @Router /grants/{id} [patch]
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

// DeleteGrant godoc
// @Summary Удалить грант
// @Tags grants
// @Param id path int true "ID гранта"
// @Success 204 "No Content"
// @Router /grants/{id} [delete]м
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

// GetGrantProjects godoc
// @Summary Проекты, финансируемые данным грантом
// @Tags grants
// @Produce json
// @Param id path int true "ID гранта"
// @Success 200 {array} models.ProjectGrantFunding
// @Router /grants/{id}/projects [get]
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

// GetGrantBudget godoc
// @Summary Состояние бюджета гранта
// @Description Возвращает общую сумму, распределенную сумму и остаток
// @Tags grants
// @Produce json
// @Param id path int true "ID гранта"
// @Success 200 {object} map[string]interface{} "Объект с полями total, allocated, remaining"
// @Router /grants/{id}/budget [get]
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

// AddProjectToGrant godoc
// @Summary Привязать проект к гранту
// @Description Выделяет определенную сумму из гранта на конкретный проект
// @Tags grants
// @Accept json
// @Produce json
// @Param id path int true "ID гранта"
// @Param funding body models.ProjectGrantFunding true "Данные финансирования"
// @Success 200 {object} models.ProjectGrantFunding
// @Router /grants/{id}/projects [post]
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

// RemoveProjectFromGrant godoc
// @Summary Удалить привязку проекта к гранту
// @Tags grants
// @Param grantId path int true "ID гранта"
// @Param fundingId path int true "ID записи финансирования"
// @Success 204 "No Content"
// @Router /grants/{grantId}/projects/{fundingId} [delete]
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

// GetProjectGrants godoc
// @Summary Гранты, привязанные к проекту
// @Description Получить список всех источников финансирования для конкретного проекта
// @Tags projects
// @Produce json
// @Param id path int true "ID проекта"
// @Success 200 {array} models.ProjectGrantInfo
// @Router /projects/{id}/grants [get]
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
