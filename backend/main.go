package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type HealthResponse struct {
	Status string `json:"status"`
	Envi   string `json:"environment"`
}

type CreateOrderReq struct {
	CustName string `json:"customer_name"`
	ProdType string `json:"product_type"`
}

func HealthChecker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{
		Status: "success",
		Envi:   os.Getenv("APP_ENV"),
	}

	json.NewEncoder(w).Encode(response)
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateOrderReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"message": "Order Created Successfully for " + req.CustName,
		"type":    req.ProdType,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
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

	http.HandleFunc("/api/ping", HealthChecker)
	http.HandleFunc("/api/v1/orders", CreateOrder)

	errs := http.ListenAndServe(":"+port, nil)
	if errs != nil {
		log.Fatalf("Server failed to start: %v", errs)
	}

	// fmt.Println("Server is Runinng on port: ", port)
	// fmt.Printf("Current environtment: %s\n", env)
}
