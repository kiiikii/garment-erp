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

	// 1. Sampling Phase
	if currentStatus == "WAITING_FOR_SAMPLE" {
		if newStatus != "SAMPLE_APPROVED" && newStatus != "SAMPLE_REVISED" && newStatus != "SAMPLE_REJECTED" {
			return errors.New("Invalid transition. Must be APPROVED, REVISED, or REJECTED.")
		}
	}

	if currentStatus == "SAMPLE_REVISED" {
		if newStatus != "WAITING_FOR_SAMPLE" && newStatus != "SAMPLE_APPROVED" {
			return errors.New("Revised sample must go back to waiting for sampling or be approved.")
		}
	}

	if currentStatus == "SAMPLE_REJECTED" {
		return errors.New("Cannot update status. This order has been sample-rejected and closed.")
	}

	// 2. LoA Phase Transitions
	if currentStatus == "SAMPLE_APPROVED" {
		if newStatus != "LOA_SIGNED" && newStatus != "LOA_REVISED" && newStatus != "LOA_REJECTED" {
			return errors.New("After sample approval, LoA must be SIGNED, REVISED, or REJECTED.")
		}
	}

	if currentStatus == "LOA_REVISED" {
		if newStatus != "LOA_SIGNED" && newStatus != "LOA_REJECTED" && newStatus != "SAMPLE_APPROVED" {
			return errors.New("Revised LoA must eventually be SIGNED, REJECTED, or sent back for approval.")
		}
	}

	if currentStatus == "LOA_REJECTED" {
		return errors.New("Cannot update status. The LoA was rejected and this order is closed.")
	}

	if currentStatus == "LOA_SIGNED" {
		if newStatus != "WAITING_FOR_MATERIALS" && newStatus != "READY_FOR_PRODUCTION" {
			return errors.New("After LoA is signed, order must either wait for materials or be marked ready for production.")
		}
	}

	if currentStatus == "WAITING_FOR_MATERIALS" {
		if newStatus != "READY_FOR_PRODUCTION" {
			return errors.New("An order waiting for materials can only transition to READY_FOR_PRODUCTION once resolved.")
		}
	}

	// 3. Production Readiness Gate
	if newStatus == "IN_PRODUCTION" {
		if currentStatus != "READY_FOR_PRODUCTION" {
			return errors.New("Can't start production. Order has not passed the Production Readiness check")
		}

		if layoutID == nil {
			return errors.New("Can't start production. PPIC hasn't assigned a sewing line layout.")
		}
	}

	// 4. Mass Production pipeline
	if currentStatus == "IN_PRODUCTION" {
		if newStatus != "CUTTING" {
			return errors.New("Order must go to CUTTING after entering production.")
		}
	}

	if currentStatus == "CUTTING" {
		if newStatus != "SEWING" {
			return errors.New("Order must go to SEWING after cutting.")
		}
	}

	if currentStatus == "SEWING" {
		if newStatus != "QC" {
			return errors.New("Order must go to QC after sewing.")
		}
	}

	if currentStatus == "QC" {
		if newStatus != "FINISHING" && newStatus != "SEWING" {
			return errors.New("After QC, order must proceed to FINISHING or go back to SEWING for rework.")
		}
	}

	if currentStatus == "FINISHING" {
		if newStatus != "READY_FOR_SHIPPING" {
			return errors.New("After finishing, order must be marked READY_FOR_SHIPPING.")
		}
	}

	return s.repo.UpdateOrderStatus(orderID, newStatus)
}

func (s *OrderService) SaveImages(orderID int, fileURL string) error {
	return s.repo.InsertOrderImage(orderID, fileURL)
}

func (s *OrderService) AssignLayout(orderID int, layoutID int) error {
	if layoutID <= 0 {
		return errors.New("Invalid layout ID provided")
	}

	return s.repo.AssignLayout(orderID, layoutID)
}

func (s *OrderService) CompleteSampling(orderID int, finishedAt string) error {
	return s.repo.CompleteSampling(orderID, finishedAt)
}

func (s *OrderService) GetByID(id int) (*OrderResponse, error) {
	return s.repo.GetOrderByID(id)
}

func (s *OrderService) SetWaitingForMaterials(orderID int, reason string) error {
	return s.repo.SetWaitingForMaterials(orderID, reason)
}

func (s *OrderService) ResumeOrder(orderID int) error {
	return s.repo.ResumeOrder(orderID)
}
