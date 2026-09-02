package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kiiikii/garment-erp/backend/internal/handlers"
	"github.com/kiiikii/garment-erp/backend/internal/repository"
	"github.com/kiiikii/garment-erp/backend/internal/services"
	_ "github.com/lib/pq"
)

func main() {

	//! Load godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	//! Read the variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//! Read APP_ENV
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}

	repo := &repository.OrderRepo{DB: db}
	svc := &services.OrderService{Repo: repo}
	app := &handlers.App{Service: svc}

	fmt.Println("Success connected to the database")

	http.HandleFunc("/api/ping", app.HealthChecker)
	http.HandleFunc("/api/v1/orders", app.OrdersRouter)
	http.HandleFunc("/api/v1/orders/status", app.UpdateOrderStatus)

	errs := http.ListenAndServe(":"+port, nil)
	if errs != nil {
		log.Fatalf("Server failed to start: %v", errs)
	}

	// fmt.Println("Server is Runinng on port: ", port)
	// fmt.Printf("Current environtment: %s\n", env)
}
