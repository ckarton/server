package models

// Response содержит данные, возвращаемые на фронт
type Response struct {
    Email       string `json:"email"`
    AccessToken string `json:"access_token"` 
    RefreshToken string `json:"refresh_token"`
    IsTeacher   bool   `json:"isTeacher"`
}

