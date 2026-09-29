package orders

import (
	"errors"
	"time"
)

type OrderService struct {
	repo *OrderRepository
}

func NewOrderService(repo *OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(o *Order) (int, error) {
	if o.TotalQuantity <= 0 {
		return 0, errors.New("Quantity must be higher than 0")
	}

	if o.ProductionType != "CMT" && o.ProductionType != "FOB" {
		return 0, errors.New("Production type must be CMT or FOB")
	}

	//! get time now
	now := time.Now()

	//! add 4 days for internal deadline
	o.InternalSampleDeadline = now.Add(4 * 24 * time.Hour)

	//! add 7 days for customer deadline
	o.CustomerSampleDeadline = now.Add(7 * 24 * time.Hour)

	return s.repo.Create(o)
}

func (s *OrderService) GetAllOrders() ([]OrderResponse, error) {
	return s.repo.GetAllWithCustomer()
}

func (s *OrderService) UpdateStatus(orderID int, newStatus string) error {
	currentStatus, layoutID, err := s.repo.GetOrderState(orderID)
	if err != nil {
		return err
	}

	//! cannot jump straight to mass product
	if newStatus == "IN_PRODUCTION" {
		if currentStatus != "LOA_SIGNED" {
			return errors.New("Can't start production. LoA hasn't been signed.")
		}

		if layoutID == nil {
			return errors.New("Can't start production. PPIC hasn't assigned")
		}
	}

	//! cannot sign Loa if sample hasn't been approved
	if newStatus == "LOA_SIGNED" {
		if currentStatus != "SAMPLE_APPROVED" {
			return errors.New("Can't sign LoA. Sample hasn't been approved")
		}
	}

	return s.repo.UpdateOrderStatus(orderID, newStatus)
}

func (s *OrderService) SaveImages(orderID int, fileURL string) error {
	return s.repo.InsertOrderImage(orderID, fileURL)
}
