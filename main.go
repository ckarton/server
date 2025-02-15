package main

import (
	"log"
	"myapp/handlers"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	// Настройка CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// Обслуживание статических файлов
	http.Handle("/", c.Handler(http.FileServer(http.Dir("./public"))))

	// API маршруты
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	// Маршруты для работы с S3
	http.HandleFunc("/upload", handlers.UploadHandler)


	// Запуск сервера
	log.Println("Server running on port 8080...")
	if err := http.ListenAndServe(":8080", c.Handler(http.DefaultServeMux)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
