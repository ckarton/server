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

	"myapp/db"
	"myapp/models"

	"github.com/google/uuid"
)

// Константы для настроек
const (
	uploadDir  = "uploads"              // Базовая директория хранения файлов
	baseURL    = "http://localhost:8080" // Базовый URL для скачивания
	maxSize    = 50 << 20                // 50MB
)

// Допустимые форматы файлов
var allowedExtensions = map[string]bool{
	".jpg":  true, ".jpeg": true, ".png": true, ".gif": true,
	".mp4":  true, ".mov":  true, ".avi": true,
	".mp3":  true, ".wav":  true,
	".pdf":  true, ".docx": true,
}

// UploadHandler загружает файлы в нужную папку
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Ограничение размера файла
	r.ParseMultipartForm(maxSize)

	// Получаем файлы из запроса
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "Файлы не загружены", http.StatusBadRequest)
		return
	}

	var uploadedFiles []string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Ошибка при открытии файла", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Проверяем допустимое расширение файла
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
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
		_, err = io.Copy(outFile, file)
		if err != nil {
			log.Printf("Ошибка сохранения файла: %v", err)
			http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
			return
		}

		// Формируем URL для доступа к файлу
		mediaURL := fmt.Sprintf("%s/uploads/%s/%s", baseURL, getFolderByType(fileType), newFileName)
		uploadedFiles = append(uploadedFiles, mediaURL)

		// Сохраняем информацию о файле в базу данных
		fileInfo := models.FileInfo{
			FileName:   fileHeader.Filename,
			FileURL:    mediaURL,
			FileType:   fileType,
			UploadedBy: r.FormValue("uploadedBy"), // предполагается, что email пользователя передается в форме
			Size:       fileHeader.Size,
		}

		_, err = db.SaveFile(fileInfo)
		if err != nil {
			log.Printf("Ошибка сохранения информации о файле: %v", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}
	}

	// Отправляем JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]string{"mediaUrls": uploadedFiles})
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
