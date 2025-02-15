package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// UploadHandler загружает файлы в нужную папку
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Ограничение размера файла (до 50MB)
	r.ParseMultipartForm(50 << 20)

	// Получаем файл из запроса
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Ошибка при получении файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Получаем тип контента (куда сохранить файл)
	fileType := r.FormValue("type")

	// Определяем папку для хранения
	baseDir := "uploads"
	var saveDir string

	switch fileType {
	case "media":
		saveDir = filepath.Join(baseDir, "media")
	case "avatars":
		saveDir = filepath.Join(baseDir, "avatars")
	case "icons":
		saveDir = filepath.Join(baseDir, "icons")
	case "courses":
		saveDir = filepath.Join(baseDir, "courses")
	default:
		saveDir = filepath.Join(baseDir, "other")
	}

	// Создаем папку, если её нет
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		os.MkdirAll(saveDir, 0755)
	}

	// Генерируем уникальное имя файла
	ext := filepath.Ext(handler.Filename) // Получаем расширение
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(saveDir, newFileName)

	// Создаем файл для сохранения
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Ошибка при создании файла", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Копируем файл в папку
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
		return
	}

	// Формируем URL для доступа к файлу
	baseURL := "http://localhost:8080"
	mediaURL := fmt.Sprintf("%s/%s", baseURL, filepath.ToSlash(filePath))

	// Отправляем JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mediaUrl": mediaURL})
}
