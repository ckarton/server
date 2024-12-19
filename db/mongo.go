package db

import (
	"context"
	"myapp/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

func init() {
	// Подключение к MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(context.Background(), clientOptions)
	collection = client.Database("myapp").Collection("users")
}

// SaveUser сохраняет пользователя в базу данных
func SaveUser(user models.User) error {
	_, err := collection.InsertOne(context.Background(), user)
	return err
}

// GetUserByEmail получает пользователя по email
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := collection.FindOne(context.Background(), map[string]string{"email": email}).Decode(&user)
	return user, err
}

// UserExists проверяет, существует ли пользователь с таким email
func UserExists(email string) bool {
	var user models.User
	err := collection.FindOne(context.Background(), map[string]string{"email": email}).Decode(&user)
	return err == nil
}
