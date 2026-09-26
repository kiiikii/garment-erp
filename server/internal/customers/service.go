package customers

import "errors"

//! create struct customer service
type CustomerService struct {
	repo *CustomerRepository
}

//! construct new service
func NewCustomerService(repo *CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

//! method create and enforce rule a customer
func (s *CustomerService) CreateCustomer(c *Customer) (int, error) {
	//! Business Rule -> Name cannot be empty
	if c.Name == "" {
		return 0, errors.New("Customer name is required")
	}

	//! Business Rule -> Phone must be at least 13 characters
	if len(c.Phone) <= 13 {
		return 0, errors.New("Phone number must be at least 13 characters")
	}

	//! if passed, tell repository to save
	return s.repo.Create(c)
}

//! get customer
func (s *CustomerService) GetAllCustomers() ([]CustomerResponse, error) {
	return s.repo.GetAll()
}
