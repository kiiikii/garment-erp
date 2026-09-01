package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kiiikii/garment-erp/backend/internal/models"
	_ "github.com/lib/pq"
)

type App struct {
	DB *sql.DB
}

func (app *App) HealthChecker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.HealthResponse{
		Status: "success",
		Envi:   os.Getenv("APP_ENV"),
	}

	json.NewEncoder(w).Encode(response)
}

func (app *App) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateOrderReq
	var newID int

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO orders (customer_name, product_type) VALUES ($1, $2) RETURNING id`
	err = app.DB.QueryRow(query, req.CustName, req.ProdType).Scan(&newID)
	if err != nil {
		http.Error(w, "Cannot insert data", http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"message":  "Order Created Successfully for " + req.CustName,
		"order_id": newID,
		"customer": req.CustName,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}

func (app *App) OrdersRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.GetOrders(w, r)
	case http.MethodPost:
		app.CreateOrder(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *App) GetOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := `SELECT id, customer_name, product_type, status FROM orders ORDER BY id DESC`
	rows, err := app.DB.Query(query)
	if err != nil {
		http.Error(w, "Database query Failed", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var orders []models.Order = []models.Order{}

	for rows.Next() {
		var o models.Order

		err := rows.Scan(&o.ID, &o.CustName, &o.ProdType, &o.Status)
		if err != nil {
			http.Error(w, "Error scanning order row", http.StatusInternalServerError)
			return
		}

		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error reading order rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

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
	app := &App{DB: db}

	fmt.Println("Success connected to the database")

	http.HandleFunc("/api/ping", app.HealthChecker)
	http.HandleFunc("/api/v1/orders", app.OrdersRouter)

	errs := http.ListenAndServe(":"+port, nil)
	if errs != nil {
		log.Fatalf("Server failed to start: %v", errs)
	}

	// fmt.Println("Server is Runinng on port: ", port)
	// fmt.Printf("Current environtment: %s\n", env)
}
