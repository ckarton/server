package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

	// Создаем путь для сохранения файла
	filePath := filepath.Join(saveDir, handler.Filename)
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

	// Отправляем ответ с путём к файлу
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Файл загружен в: %s", filePath)
}
