package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	uploadDir = "uploads"
	baseURL   = "http://localhost:8080"
	maxSize   = 50 << 20 // 50MB
)

var allowedExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".mp4": true, ".mov": true, ".avi": true,
	".mp3": true, ".wav": true,
	".pdf": true, ".docx": true,
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Ограничение размера файла
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		http.Error(w, "Файл слишком большой", http.StatusRequestEntityTooLarge)
		return
	}
	// Получаем файл из запроса
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Ошибка при получении файла", http.StatusBadRequest)
		return
	}
	defer file.Close()
	// Проверяем допустимое расширение файла
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if _, valid := allowedExtensions[ext]; !valid {
		http.Error(w, "Недопустимый формат файла", http.StatusUnsupportedMediaType)
		return
	}

	// Получаем тип контента (куда сохранить файл)
	fileType := r.FormValue("type")
	// Определяем папку для хранения
	saveDir := filepath.Join(uploadDir, getFolderByType(fileType))
	// Создаем папку, если её нет
	err = os.MkdirAll(saveDir, 0755)
	if err != nil {
		log.Printf("Ошибка создания папки: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	// Генерируем уникальное имя файла
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(saveDir, newFileName)
	// Создаем файл для сохранения
	outFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("Ошибка создания файла: %v", err)
		http.Error(w, "Ошибка при создании файла", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()
	// Копируем файл в папку
	size, err := io.Copy(outFile, file)
	if err != nil {
		log.Printf("Ошибка сохранения файла: %v", err)
		http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
		return
	}

	clientIP := r.RemoteAddr
	log.Printf("Файл загружен: %s (Оригинальное имя: %s, Размер: %d байт, Тип: %s, IP: %s)", newFileName, handler.Filename, size, fileType, clientIP)

	mediaURL := fmt.Sprintf("%s/uploads/%s/%s", baseURL, getFolderByType(fileType), newFileName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mediaUrl": mediaURL})
}
// getFolderByType возвращает соответствующую папку по типу контента
func getFolderByType(fileType string) string {
	switch fileType {
	case "media":
		return "media"
	case "avatars":
		return "avatars"
	case "icons":
		return "icons"
	case "courses":
		return "courses"
	default:
		return "other"
	}
}
