package main

import (
	"em-test/cmd/config"
	"em-test/cmd/internal/db"
	"em-test/cmd/internal/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "em-test/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title EM-test
// @version 1.0
// @description REST API для управления подписками
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.Load()
	database := db.ConnectPostgres(cfg.DSN)

	//Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		sqlDB, err := database.DB()
		if err == nil {
			sqlDB.Close()
		}
		log.Printf("[%v] Subscription server stopped: DB-connections closed.\n", time.Now().Format("2006-01-02 15:04:05"))
		os.Exit(0)
	}()

	//Creting hadnler with embedded service and repo
	subHandler := handler.CreateHandler(database)
	r := chi.NewRouter()

	//HTTP-handlers: service and swagger
	r.Post("/subscriptions", subHandler.Create)
	r.Get("/subscriptions", subHandler.GetList)
	r.Get("/subscriptions/{sid}", subHandler.GetBySID)
	r.Delete("/subscriptions/{sid}", subHandler.Delete)
	r.Put("/subscriptions/{sid}", subHandler.UpdateBySID)

	r.Get("/subscriptions/report", subHandler.Report)
	//GET  /subscriptions/report?period=05-2024&uid=42&provider=YoutubePremium

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	//Starting server
	log.Printf("Server running on http://localhost%s", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
