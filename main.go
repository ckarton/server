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
		AllowedOrigins: []string{"*"}, // Разрешить запросы с любого источника (можно ограничить для продакшн)
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// Обслуживание статических файлов с CORS
	http.Handle("/", c.Handler(http.FileServer(http.Dir("./public"))))

	// Настройка API маршрутов с CORS
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	// Запуск сервера
	log.Println("Server running on port 8080...")
	if err := http.ListenAndServe(":8080", c.Handler(http.DefaultServeMux)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
