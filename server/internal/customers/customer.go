package customers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ! defining struct
type Customer struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type CustomerResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

// ! handler
func CreateCustomerHandler(db *sql.DB) http.HandlerFunc {
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

func GetCustomerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//! enforce http method GET
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed. Use GET.", http.StatusMethodNotAllowed)
			return
		}

		//! ask postgres for all customer
		rows, err := db.Query("SELECT id, name, phone, address FROM customers")
		if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}

		//! close the rows
		defer rows.Close()

		//! create an empty slice
		var customers []CustomerResponse

		//! Loop through the row
		for rows.Next() {
			var c CustomerResponse

			//! scan the raw sql
			err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Address)
			if err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				return
			}

			//! add the struct to list
			customers = append(customers, c)
		}

		//! telling client sending back JSON data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		//! translate go list to JSON
		json.NewEncoder(w).Encode(customers)
	}
}
