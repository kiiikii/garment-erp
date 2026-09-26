package orders

import (
	"database/sql"
	"errors"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// ! create function insert data
func (r *OrderRepository) Create(o *Order) (int, error) {
	var newID int
	sqlStatement := `INSERT INTO orders (customer_id, quantity, production_type) VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRow(sqlStatement, o.CustomerID, o.Quantity, o.ProductionType).Scan(&newID)
	if err != nil {
		return 0, errors.New("Failed to insert data, is the Customer Exist")
	}

	return newID, nil
}

// ! JOIN query
func (r *OrderRepository) GetAllWithCustomer() ([]OrderResponse, error) {
	sqlStatement := `SELECT o.id, o.quantity, o.production_type, o.created_at, c.id, c.name FROM orders o JOIN customers c ON o.customer_id = c.id`

	rows, err := r.db.Query(sqlStatement)
	if err != nil {
		return nil, errors.New("Failed to query orders")
	}

	defer rows.Close()

	var orderList []OrderResponse

	for rows.Next() {
		var res OrderResponse
		err := rows.Scan(&res.ID, &res.Quantity, &res.ProductionType, &res.CreatedAt, &res.Customer.ID, &res.Customer.Name)
		if err != nil {
			return nil, errors.New("Failed to scan order row")
		}
		orderList = append(orderList, res)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("Database connection dies during the loop")
	}

	return orderList, nil
}
