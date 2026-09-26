package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/kiiikii/garment-erp/server/internal/customers"
	"github.com/kiiikii/garment-erp/server/internal/orders"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Garment API is running!")
}

// ! CORS Middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//! telling the browser where's NextJS
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

		//! specify permitted http mehtods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		//! allowing Next.JS send JSON data
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	//! Load the env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found. Relying on system env")
	}

	//! get the database URL
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("CRITICAL: DB_URL env variable is not set")
	}

	//! connect the PostgreSQL
	fmt.Println("Connecting the database...")
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("Failed to open database connection: ", err)
	}

	//! ensure the database connection close
	defer db.Close()

	//! ping the database
	err = db.Ping()
	if err != nil {
		log.Println("WARNING: Couldn't ping database. Is PostgreSQL running?", err)
	} else {
		fmt.Println("SUCCESS: Connected to PostgreSQL")
	}

	//! starting http server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//! repository
	customerRepo := customers.NewCustomeRepository(db)
	orderRepo := orders.NewOrderRepository(db)

	//! service
	customerService := customers.NewCustomerService(customerRepo)
	orderService := orders.NewOrderService(orderRepo)

	//! handler
	customerHandler := customers.NewCustomerHandler(customerService)
	orderHandler := orders.NewOrderHandler(orderService)

	http.HandleFunc("/health", healthCheckHandler)

	http.HandleFunc("/customers", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			customerHandler.CreateCustomer(w, r)
		} else if r.Method == http.MethodGet {
			customerHandler.GetAllCustomers(w, r)
		} else {
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/orders", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			orderHandler.CreateOrder(w, r)
		} else if r.Method == http.MethodGet {
			orderHandler.GetAllOrders(w, r)
		} else {
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	fmt.Printf("Starting Garment API server on port %s...\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server Crashed: ", err)
	}
}
