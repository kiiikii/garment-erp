package main

import "fmt"

type Order struct {
	ID           int
	CustomerName string
	IsWaiting    bool
}

type Describable interface {
	Describe() string
}

func (o Order) Describe() string {
	if o.IsWaiting {
		return fmt.Sprintf("Order for, %s is waiting.", o.CustomerName)
	} else {
		return fmt.Sprintf("Order for, %s is not waiting.", o.CustomerName)
	}
}

func Description(d Describable) {
	fmt.Println(d.Describe())
}

func validateOrder(o Order) error {
	if o.CustomerName == "" {
		return fmt.Errorf("customer name is required")
	}

	return nil
}

func main() {
	order := Order{
		ID:           123,
		CustomerName: "Acme Corp",
		IsWaiting:    false,
	}

	order1 := Order{
		ID:           122,
		CustomerName: "",
		IsWaiting:    true,
	}

	err := validateOrder(order1)
	if err != nil {
		fmt.Println("validation Failed: ", err.Error())
	} else {
		fmt.Println("validation Success")
	}

	fmt.Println(
		"Order = ", order.ID, ", customer name = ", order.CustomerName, ", is waiting = ", order.IsWaiting,
	)

	Description(order)
}
