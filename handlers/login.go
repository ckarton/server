package handlers

import (
	"encoding/json"
	"myapp/db"
	"myapp/models"
	"myapp/utils"
	"net/http"
)

// LoginHandler для логина пользователя
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var input models.User
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Получаем пользователя из базы
    user, err := db.GetUserByEmail(input.Email)
    if err != nil {
        http.Error(w, "User not found", http.StatusUnauthorized)
        return
    }

    // Проверяем пароль
    if !utils.CheckPasswordHash(input.Password, user.Password) {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    // Генерация токенов
    accessToken, refreshToken, err := GenerateTokens(user.Email)
    if err != nil {
        http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
        return
    }

    // Ответ на фронт
    response := models.Response{
        Email:       user.Email,
        AccessToken: accessToken,
        RefreshToken: refreshToken,
        IsTeacher:   user.IsTeacher,
    }

    json.NewEncoder(w).Encode(response)
}
