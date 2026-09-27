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
	//! transaction begin
	tx, err := r.db.Begin()
	if err != nil {
		return 0, errors.New("Failed to start database transaction")
	}

	//! define the defer rollback
	defer tx.Rollback()

	var newOrderID int

	//! insert main order
	orderSQL := `INSERT INTO orders (customer_id, total_quantity, production_type) VALUES ($1, $2, $3) RETURNING id`

	err = tx.QueryRow(orderSQL, o.CustomerID, o.TotalQuantity, o.ProductionType).Scan(&newOrderID)
	if err != nil {
		return 0, errors.New("Failed to insert main order")
	}

	//! loop and insert
	sizeSQL := `INSERT INTO order_sizes (order_id, size_label, quantity) VALUES ($1, $2, $3)`

	for _, size := range o.Sizes {
		//! use tx.Exec t insert each size, linking newOrderID
		_, err := tx.Exec(sizeSQL, newOrderID, size.SizeLabel, size.Quantity)
		if err != nil {
			//! if even one size fail, return error, the defer tx.Rollback will instant erase main data
			return 0, errors.New("Failed to insert order size: " + size.SizeLabel)
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, errors.New("failed to commit transaction")
	}

	return newOrderID, nil
}

// ! JOIN query
func (r *OrderRepository) GetAllWithCustomer() ([]OrderResponse, error) {
	sqlStatement := `SELECT o.id, o.total_quantity, o.production_type, o.created_at, c.id, c.name FROM orders o JOIN customers c ON o.customer_id = c.id`

	rows, err := r.db.Query(sqlStatement)
	if err != nil {
		return nil, errors.New("Failed to query orders")
	}

	defer rows.Close()

	var orderList []OrderResponse

	for rows.Next() {
		var res OrderResponse

		//! update scan
		err := rows.Scan(&res.ID, &res.TotalQuantity, &res.ProductionType, &res.CreatedAt, &res.Customer.ID, &res.Customer.Name)
		if err != nil {
			return nil, errors.New("Failed to scan order row")
		}

		//! fetch the size for specific order id
		sizeSQL := `SELECT size_label, quantity FROM order_sizes WHERE order_id = $1`
		sizeRows, err := r.db.Query(sizeSQL, res.ID)
		if err != nil {
			return nil, errors.New("Failed to query sizes")
		}

		//! loop through size and add them to res.Sizes slices
		for sizeRows.Next() {
			var s OrderSize
			if err := sizeRows.Scan(&s.SizeLabel, &s.Quantity); err != nil {
				sizeRows.Close()
				return nil, errors.New("Failed to scan size")
			}
			res.Sizes = append(res.Sizes, s)
		}
		sizeRows.Close()

		//! append fully built order
		orderList = append(orderList, res)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("Database connection dies during the loop")
	}

	return orderList, nil
}
