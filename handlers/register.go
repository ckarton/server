package handlers

import (
	"encoding/json"
	"myapp/db"
	"myapp/models"
	"myapp/utils"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Хеш пароля
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Генерация токенов
	accessToken, refreshToken, err := GenerateTokens(user.Email)
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	// Записываем пользователя в Mongo
	user.Password = hashedPassword
	user.AccessToken = accessToken
	user.RefreshToken = refreshToken
	err = db.SaveUser(user)
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	// Ответ на фронт
	response := models.Response{
		Email:     user.Email,
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		IsTeacher: user.IsTeacher,
	}

	json.NewEncoder(w).Encode(response)
}
