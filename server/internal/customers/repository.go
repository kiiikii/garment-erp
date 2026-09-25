package customers

import (
	"database/sql"
	"errors"
)

// ! create struct repository
type CustomerRepository struct {
	db *sql.DB
}

// ! create func construct to create new repository
func NewCustomeRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// ! method insert customer
func (r *CustomerRepository) Create(c *Customer) (int, error) {
	var newID int
	sqlStatement := `INSERT INTO customers (name, phone, address) VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRow(sqlStatement, c.Name, c.Phone, c.Address).Scan(&newID)
	if err != nil {
		return 0, errors.New("Failed to insert customer into database")
	}

	return newID, nil
}

func (r *CustomerRepository) GetAll() ([]CustomerResponse, error) {
	rows, err := r.db.Query("SELECT id, name, phone, address FROM customers")
	if err != nil {
		return nil, errors.New("failed to query customers")
	}

	defer rows.Close()

	var customerList []CustomerResponse

	for rows.Next() {
		var c CustomerResponse
		err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Address)
		if err != nil {
			return nil, errors.New("failed to scan customer row")
		}
		customerList = append(customerList, c)
	}

	return customerList, nil
}
