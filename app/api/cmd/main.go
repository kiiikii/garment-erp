package main

import (
	"fmt"
	"garment-erp-api/internal/order"
)

func main() {
	validOrder := order.Order{
		ID:           123,
		CustomerName: "Acme Corp",
		IsWaiting:    false,
	}

	invalidOrder := order.Order{
		ID:           122,
		CustomerName: "",
		IsWaiting:    true,
	}

	err := order.ValidateOrder(invalidOrder)
	if err != nil {
		fmt.Println("validation Failed: ", err.Error())
	} else {
		fmt.Println("validation Success")
	}

	fmt.Println(
		"Order = ", validOrder.ID, ", customer name = ", validOrder.CustomerName, ", is waiting = ", validOrder.IsWaiting,
	)

	order.Description(validOrder)
}
