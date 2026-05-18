package handlers

import (
	"encoding/json"
	"net/http"
	"project-MVP/services"
)

// GetActivitiesHandler godoc
// @Summary Получить ленту активности
// @Description Возвращает список последних действий (создание задач, комментариев, загрузка файлов) во всех проектах, где состоит текущий пользователь
// @Tags activity
// @Produce json
// @Success 200 {array} models.Activity "Список событий активности"
// @Failure 401 {string} string "Пользователь не авторизован"
// @Failure 500 {string} string "Ошибка при получении данных из базы"
// @Router /activities [get]
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
