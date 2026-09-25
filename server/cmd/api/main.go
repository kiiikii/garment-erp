package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

// ! defining struct
type Customer struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

// ! handler
func createCustomerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//! enforce using HTTP POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed. Use POST", http.StatusMethodNotAllowed)
			return
		}

		//! decode incoming JSON
		var c Customer
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		//! inserting into postgres
		var newID int
		sqlStatement := `INSERT INTO customers (name, phone, address) VALUES ($1, $2, $3) RETURNING id`

		//! use db.QueryRow
		err = db.QueryRow(sqlStatement, c.Name, c.Phone, c.Address).Scan(&newID)
		if err != nil {
			log.Println("Database Error:", err)
			http.Error(w, "Failed to save customer to database", http.StatusInternalServerError)
			return
		}

		//! sending success response back
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Successfully created Customer #%d: %s\n", newID, c.Name)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Garment API is running!")
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

	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/customers", createCustomerHandler(db))

	fmt.Printf("Starting Garment API server on port %s...\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server Crashed: ", err)
	}
}
