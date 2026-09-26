package orders

import "errors"

type OrderService struct {
	repo *OrderRepository
}

func NewOrderService(repo *OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(o *Order) (int, error) {
	if o.Quantity <= 0 {
		return 0, errors.New("Quantity must be higher than 0")
	}

	if o.ProductionType != "CMT" && o.ProductionType != "FOB" {
		return 0, errors.New("Production type must be CMT or FOB")
	}

	return s.repo.Create(o)
}

func (s *OrderService) GetAllOrders() ([]OrderResponse, error) {
	return s.repo.GetAllWithCustomer()
}
