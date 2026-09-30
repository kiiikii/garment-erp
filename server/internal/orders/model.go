package orders

import "time"

//! size struct
type OrderSize struct {
	SizeLabel string `json:"size_label"`
	Quantity  int    `json:"quantity"`
}

//! status update struct
type StatusUpdateRequest struct {
	CurrentStatus string `json:"current_status"`
	NewStatus     string `json:"new_status"`
}

//! create order struct
type Order struct {
	CustomerID             int         `json:"customer_id"`
	TotalQuantity          int         `json:"total_quantity"`
	ProductionType         string      `json:"production_type"`
	Status                 string      `json:"status"`
	InternalSampleDeadline time.Time   `json:"-"`
	CustomerSampleDeadline time.Time   `json:"-"`
	Sizes                  []OrderSize `json:"sizes"`
}

//! create response struct
type OrderResponse struct {
	ID                     int        `json:"id"`
	TotalQuantity          int        `json:"total_quantity"`
	ProductionType         string     `json:"production_type"`
	Status                 string     `json:"status"`
	CreatedAt              time.Time  `json:"created_at"`
	InternalSampleDeadline time.Time  `json:"internal_sample_deadline"`
	CustomerSampleDeadline time.Time  `json:"customer_sample_deadline"`
	ActualSampleFinishedAt *time.Time `json:"actual_sample_finished_at"`
	WaitingReason          *string    `json:"waiting_reason"`
	Customer               struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
	} `json:"customer"`
	Sizes  []OrderSize `json:"sizes"`
	Images []struct {
		ID       int    `json:"id"`
		ImageURL string `json:"image_url"`
	} `json:"images"`
	LayoutID *int `json:"layout_id"`
}
