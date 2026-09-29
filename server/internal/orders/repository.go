package orders

import (
	"database/sql"
	"errors"
	"fmt"
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
	orderSQL := `INSERT INTO orders (customer_id, total_quantity, production_type, internal_sample_deadline, customer_sample_deadline) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err = tx.QueryRow(orderSQL, o.CustomerID, o.TotalQuantity, o.ProductionType, o.InternalSampleDeadline, o.CustomerSampleDeadline).Scan(&newOrderID)
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
	sqlStatement := `SELECT o.id, o.total_quantity, o.production_type, o.created_at, o.status, o.internal_sample_deadline, o.customer_sample_deadline, o.actual_sample_finished_at, c.id, c.name FROM orders o JOIN customers c ON o.customer_id = c.id`

	rows, err := r.db.Query(sqlStatement)
	if err != nil {
		fmt.Println("DB Query Error:", err)
		return nil, errors.New("Failed to query orders")
	}

	defer rows.Close()

	var orderList []OrderResponse

	for rows.Next() {
		var res OrderResponse
		var safeActualFinished sql.NullTime

		//! update scan
		err := rows.Scan(&res.ID, &res.TotalQuantity, &res.ProductionType, &res.CreatedAt, &res.Status, &res.InternalSampleDeadline, &res.CustomerSampleDeadline, &safeActualFinished, &res.Customer.ID, &res.Customer.Name)
		if err != nil {
			return nil, errors.New("Failed to scan order row")
		}
		if safeActualFinished.Valid {
			res.ActualSampleFinishedAt = &safeActualFinished.Time
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

		if err := sizeRows.Err(); err != nil {
			return nil, errors.New("Database Connection dies during the loop")
		}

		//! append fully built order
		orderList = append(orderList, res)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("Database connection dies during the loop")
	}

	return orderList, nil
}

// ! update status order
func (r *OrderRepository) UpdateOrderStatus(orderID int, newStatus string) error {
	sqlStatement := `UPDATE orders SET status = $1 WHERE id = $2`

	_, err := r.db.Exec(sqlStatement, newStatus, orderID)
	if err != nil {
		return errors.New("failed to update status in database")
	}

	return nil
}

// ! insert order images
func (r *OrderRepository) InsertOrderImage(orderID int, fileURL string) error {
	sqlStatement := `INSERT INTO order_images (order_id, file_url) VALUES ($1, $2)`
	_, err := r.db.Exec(sqlStatement, orderID, fileURL)
	if err != nil {
		return errors.New("Failed to saved image record to database")
	}

	return nil
}
