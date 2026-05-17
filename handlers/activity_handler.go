package handlers

import (
	"encoding/json"
	"net/http"
	"project-MVP/services"
)

func GetActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Проверяем авторизацию и получаем ID пользователя
	userID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	// 2. Вызываем сервис для получения списка активностей
	list, err := services.GetUserActivities(userID)
	if err != nil {
		http.Error(w, "Ошибка получения ленты: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Отправляем JSON клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
