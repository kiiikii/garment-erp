package orders

import "time"

//! size struct
type OrderSize struct {
	SizeLabel string `json:"size_label"`
	Quantity  int    `json:"quantity"`
}

//! create order struct
type Order struct {
	CustomerID     int         `json:"customer_id"`
	TotalQuantity  int         `json:"total_quantity"`
	ProductionType string      `json:"production_type"`
	Status         string      `json:"status"`
	Sizes          []OrderSize `json:"sizes"`
}

//! create response struct
type OrderResponse struct {
	ID             int       `json:"id"`
	TotalQuantity  int       `json:"total_quantity"`
	ProductionType string    `json:"production_type"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	Customer       struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"customer"`
	Sizes []OrderSize `json:"sizes"`
}
