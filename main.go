package main

import (
	"log"
	"net/http"
	"time"

	"architecture/internal"
)

const inventoryPath = "data/inventory.csv"

func main() {
	repo := internal.NewPartRepository()
	service := internal.NewPartService(repo)
	handler := internal.NewHandler(service)

	if err := repo.LoadFromCSV(inventoryPath); err != nil {
		log.Printf("Не удалось загрузить данные: %v", err)
	}

	log.Println("Сервер запущен на http://localhost:8080")

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
