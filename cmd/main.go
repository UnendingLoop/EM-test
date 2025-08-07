package main

import (
	"em-test/cmd/internal/db"
	"em-test/cmd/internal/handler"
	"em-test/cmd/internal/repository"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title EM-test
// @version 1.0
// @description REST API для управления подписками
// @host localhost:8080
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set in env")
	}
	db := db.ConnectPostgres(dsn)

	subRepo := repository.SubscriptionRepo{DB: db}
	subHandler := handler.SubscriptionHandler{Repo: subRepo}

	r := chi.NewRouter()

	r.Post("/subscription/create", subHandler.Create)
	r.Get("/subscription/{id}", subHandler.GetByID)
	r.Get("/subscription/list", subHandler.GetList)
	r.Delete("/subscription/delete/{id}", subHandler.Delete)
	r.Put("/subscription/update/{id}", subHandler.UpdateByID)

	r.Get("/search", subHandler.Search)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	port := os.Getenv("SUBSCRIPTION_PORT")
	if port == "" {
		log.Fatal("PORT is not set in env")
	}
	fmt.Printf("Server running on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}
