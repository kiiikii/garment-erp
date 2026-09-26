package orders

import "time"

//! create order struct
type Order struct {
	CustomerID     int    `json:"customer_id"`
	Quantity       int    `json:"quantity"`
	ProductionType string `json:"production_type"`
}

//! create response struct
type OrderResponse struct {
	ID             int       `json:"id"`
	Quantity       int       `json:"quantity"`
	ProductionType string    `json:"production_type"`
	CreatedAt      time.Time `json:"created_at"`
	Customer       struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
}
