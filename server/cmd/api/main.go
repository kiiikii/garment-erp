package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

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

	fmt.Printf("Starting Garment API server on port %s...\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server Crashed: ", err)
	}
}
