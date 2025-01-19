package models

// Структура пользователя
type User struct {
    Email       string `json:"email" bson:"email"`
    Password    string `json:"password" bson:"password"`
    AccessToken string `json:"access_token" bson:"access_token"` 
    RefreshToken string `json:"refresh_token" bson:"refresh_token"`
    IsTeacher   bool   `json:"isTeacher" bson:"isTeacher"`
}
