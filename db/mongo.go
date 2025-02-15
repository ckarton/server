package db

import (
	"context"
	"myapp/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var usersCollection *mongo.Collection
var filesCollection *mongo.Collection

func init() {
	// Подключение к MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(context.Background(), clientOptions)
	usersCollection = client.Database("myapp").Collection("users")
	filesCollection = client.Database("myapp").Collection("files")
}

// сохраняет пользователя в базу данных
func SaveUser(user models.User) error {
	_, err := usersCollection.InsertOne(context.Background(), user)
	return err
}

// получает пользователя по email
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := usersCollection.FindOne(context.Background(), map[string]string{"email": email}).Decode(&user)
	return user, err
}

// проверяет существует ли пользователь с таким email
func UserExists(email string) bool {
	var user models.User
	err := usersCollection.FindOne(context.Background(), map[string]string{"email": email}).Decode(&user)
	return err == nil
}

// SaveFile сохраняет информацию о файле в базу данных
func SaveFile(file models.FileInfo) (string, error) {
	if file.UploadedAt.IsZero() {
		file.UploadedAt = time.Now()
	}

	// Генерация уникального UUID для файла
	file.ID = uuid.New().String()

	_, err := filesCollection.InsertOne(context.Background(), file)
	if err != nil {
		return "", err
	}

	return file.ID, nil
}

// GetFileByID получает информацию о файле по его ID
func GetFileByID(fileID string) (models.FileInfo, error) {
	var file models.FileInfo
	err := filesCollection.FindOne(context.Background(), bson.M{"_id": fileID}).Decode(&file)
	return file, err
}

// GetFilesByUser получает все файлы, загруженные пользователем
func GetFilesByUser(userEmail string) ([]models.FileInfo, error) {
	cursor, err := filesCollection.Find(context.Background(), bson.M{"uploadedBy": userEmail})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var files []models.FileInfo
	if err = cursor.All(context.Background(), &files); err != nil {
		return nil, err
	}
	return files, nil
}

// DeleteFile удаляет информацию о файле из базы данных
func DeleteFile(fileID string) error {
	_, err := filesCollection.DeleteOne(context.Background(), bson.M{"_id": fileID})
	return err
}
