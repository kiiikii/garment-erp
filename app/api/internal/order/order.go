package order

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

func ValidateOrder(o Order) error {
	if o.CustomerName == "" {
		return fmt.Errorf("customer name is required")
	}

	return nil
}
