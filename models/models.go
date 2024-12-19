package models

// Response содержит данные, возвращаемые на фронт
type Response struct {
    Email     string `json:"email"`
    Token     string `json:"token"`
    IsTeacher bool   `json:"isTeacher"` // добавляем булево значение
}
