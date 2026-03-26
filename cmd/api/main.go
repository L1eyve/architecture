package main

import (
	"golang-arch/internal/handler"
	"golang-arch/internal/repository"
	"golang-arch/internal/service"
	"log"
	"net/http"
	"time"
)

const inventoryPath = "data/inventory.csv"

func main() {
	repo := repository.NewPartRepository()
	service := service.NewPartService(repo)
	h := handler.NewHandler(service)

	if err := repo.LoadFromCSV(inventoryPath); err != nil {
		log.Printf("Не удалось загрузить данные: %v", err)
	}

	log.Println("Сервер запущен на http://localhost:8080")

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
