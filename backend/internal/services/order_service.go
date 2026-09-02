package services

import (
	"errors"

	"github.com/kiiikii/garment-erp/backend/internal/models"
	"github.com/kiiikii/garment-erp/backend/internal/repository"
)

type OrderService struct {
	Repo *repository.OrderRepo
}

func (s *OrderService) CreateOrder(req models.CreateOrderReq) (int, error) {
	//! Bussiness Rule
	if req.ProdType != "CMT" && req.ProdType != "FOB" {
		return 0, errors.New("Invalid product type: must be CMT or FOB")
	}

	//! If valid, pass to database
	return s.Repo.InsertOrder(req)
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.Repo.GetAllOrders()
}
