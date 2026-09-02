package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/kiiikii/garment-erp/backend/internal/models"
	_ "github.com/lib/pq"
)

type App struct {
	DB *sql.DB
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
