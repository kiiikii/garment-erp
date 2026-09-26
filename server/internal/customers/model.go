package customers

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
