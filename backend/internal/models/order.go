package models

type Order struct {
	ID       int    `json:"id"`
	CustName string `json:"customer_name"`
	ProdType string `json:"product_type"`
	Status   string `json:"status"`
}

type CreateOrderReq struct {
	CustName string `json:"customer_name"`
	ProdType string `json:"product_type"`
}
