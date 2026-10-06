package customers

import (
	"errors"
	"strings"
)

// ! create struct customer service
type CustomerService struct {
	repo *CustomerRepository
}

// ! construct new service
func NewCustomerService(repo *CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

// ! method create and enforce rule a customer
func (s *CustomerService) CreateCustomer(c *Customer) (int, error) {
	//! Business Rule -> Name cannot be empty
	if strings.TrimSpace(c.Name) == "" {
		return 0, errors.New("Customer name is required")
	}

	//! Business Rule -> Phone must be at least 15 characters
	if len(strings.TrimSpace(c.Phone)) <= 10 {
		return 0, errors.New("Phone number must be at least 10 characters")
	}

	if c.Email == "" || (!strings.Contains(c.Email, "@")) {
		return 0, errors.New("Email cannot be empty & must contain a '@' symbol")
	}

	//! if passed, tell repository to save
	return s.repo.Create(c)
}

// ! get customer
func (s *CustomerService) GetAllCustomers() ([]CustomerResponse, error) {
	return s.repo.GetAll()
}
