package models

// Структура пользователя
type User struct {
    Email    string `json:"email" bson:"email"`
    Password string `json:"password" bson:"password"`
    Token    string `json:"token" bson:"token"` // поле для токена
    IsTeacher bool  `json:"isTeacher" bson:"isTeacher"`
}
