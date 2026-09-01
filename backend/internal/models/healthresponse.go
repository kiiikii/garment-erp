package models

type HealthResponse struct {
	Status string `json:"status"`
	Envi   string `json:"environment"`
}
