package notifications

import "time"

type Notification struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	AlertType string    `json:"alert_type"`
	SentAt    time.Time `json:"sent_at"`
}

type DelayedOrder struct {
	OrderID       int
	Status        string
	Deadline      time.Time
	CustomerEmail string
	CustomerPhone string
}
