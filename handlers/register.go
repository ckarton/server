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

	// Генерация токена
	token, err := GenerateToken(user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Записываем пользователя в Mongo
	user.Password = hashedPassword
	user.Token = token
	err = db.SaveUser(user)
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	// Ответ на фронт
	response := models.Response{
		Email:     user.Email,
		Token:     token,
		IsTeacher: user.IsTeacher,
	}

	json.NewEncoder(w).Encode(response)
}
