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
	sqlStatement := `
		SELECT o.id, o.total_quantity, o.production_type, o.created_at, 
		o.status, o.internal_sample_deadline, o.customer_sample_deadline, 
		o.actual_sample_finished_at, o.waiting_reason, o.layout_id,
		o.internal_production_deadline, o.customer_production_deadline,
		o.actual_production_started_at, o.actual_production_completed_at, c.id, c.name 
		FROM orders o JOIN customers c ON o.customer_id = c.id ORDER BY o.id::integer ASC`

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
		var safeInternalProductionDeadline sql.NullTime
		var safeCustomerProductionDeadline sql.NullTime
		var safeActualProductionStartedAt sql.NullTime
		var safeActualProductionCompletedAt sql.NullTime
		var safeWaitingReason sql.NullString
		var safeLayoutID sql.NullInt64

		//! update scan
		err := rows.Scan(
			&res.ID, &res.TotalQuantity, &res.ProductionType,
			&res.CreatedAt, &res.Status, &res.InternalSampleDeadline,
			&res.CustomerSampleDeadline, &safeActualFinished, &safeWaitingReason,
			&safeLayoutID, &safeInternalProductionDeadline, &safeCustomerProductionDeadline,
			&safeActualProductionStartedAt, &safeActualProductionCompletedAt,
			&res.Customer.ID, &res.Customer.Name)
		if err != nil {
			return nil, errors.New("Failed to scan order row")
		}
		if safeActualFinished.Valid {
			res.ActualSampleFinishedAt = &safeActualFinished.Time
		}
		if safeInternalProductionDeadline.Valid {
			res.InternalProductionDeadline = &safeInternalProductionDeadline.Time
		}
		if safeCustomerProductionDeadline.Valid {
			res.CustomerProductionDeadline = &safeCustomerProductionDeadline.Time
		}
		if safeActualProductionStartedAt.Valid {
			res.ActualProductionStartedAt = &safeActualProductionStartedAt.Time
		}
		if safeActualProductionCompletedAt.Valid {
			res.ActualProductionCompletedAt = &safeActualProductionCompletedAt.Time
		}
		if safeWaitingReason.Valid {
			res.WaitingReason = &safeWaitingReason.String
		}
		if safeLayoutID.Valid {
			val := int(safeLayoutID.Int64)
			res.LayoutID = &val
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
	var sqlStatement string

	switch newStatus {
	case "WAITING_FOR_SAMPLE":
		sqlStatement = `UPDATE orders SET status = $1, actual_sample_finished_at = NULL WHERE id = $2`
	case "IN_PRODUCTION":
		sqlStatement = `UPDATE orders SET status = $1, actual_production_started_at = CURRENT_TIMESTAMP,
										internal_production_deadline = CURRENT_TIMESTAMP + INTERVAL '9 days',
										customer_production_deadline = CURRENT_TIMESTAMP + INTERVAL '12 days' WHERE id = $2`
	case "READY_FOR_SHIPPING":
		sqlStatement = `UPDATE orders SET status = $1, actual_production_completed_at = CURRENT_TIMESTAMP WHERE id = $2`
	default:
		sqlStatement = `UPDATE orders SET status = $1 WHERE id = $2`
	}

	_, err := r.db.Exec(sqlStatement, newStatus, orderID)
	return err
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

// ! get order state
func (r *OrderRepository) GetOrderState(orderID int) (string, *int, error) {
	var currentStatus string
	var safeLayoutID sql.NullInt64

	query := `SELECT status, layout_id FROM orders WHERE id = $1`
	err := r.db.QueryRow(query, orderID).Scan(&currentStatus, &safeLayoutID)
	if err != nil {
		fmt.Println("GetOrderState DB Error:", err)
		return "", nil, errors.New("Failed to find order state")
	}

	var layoutID *int
	if safeLayoutID.Valid {
		id := int(safeLayoutID.Int64)
		layoutID = &id
	}

	return currentStatus, layoutID, nil
}

// ! assigning layout
func (r *OrderRepository) AssignLayout(orderID int, layoutID int) error {
	query := `UPDATE orders SET layout_id = $1 WHERE id = $2 AND status = 'READY_FOR_PRODUCTION'`
	result, err := r.db.Exec(query, layoutID, orderID)
	if err != nil {
		return errors.New("Failed to execute layout assignment")
	}

	rowAffected, err := result.RowsAffected()
	if err != nil || rowAffected == 0 {
		return errors.New("Order not found or layout update failed")
	}

	return nil
}

// ! sampling complete
func (r *OrderRepository) CompleteSampling(orderID int, finishedAt string) error {
	query := `UPDATE orders SET actual_sample_finished_at = $1 WHERE id = $2`
	result, err := r.db.Exec(query, finishedAt, orderID)
	if err != nil {
		return errors.New("Failed to record actual sample completion")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("Order not found")
	}

	return nil
}

// ! get order by id
func (r *OrderRepository) GetOrderByID(id int) (*OrderResponse, error) {
	query :=
		`SELECT o.id, o.customer_id, c.name, c.phone, c.address,
		o.total_quantity, o.production_type, o.status, o.created_at,
		o.internal_sample_deadline, o.customer_sample_deadline,
		o.actual_sample_finished_at, o.layout_id, o.waiting_reason, o.internal_production_deadline, 
		o.customer_production_deadline, o.actual_production_started_at, o.actual_production_completed_at 
		FROM orders o JOIN customers c ON o.customer_id = c.id WHERE o.id = $1`

	var res OrderResponse
	var safeActualFinished sql.NullTime
	var safeInternalProductionDeadline sql.NullTime
	var safeCustomerProductionDeadline sql.NullTime
	var safeActualProductionStartedAt sql.NullTime
	var safeActualProductionCompletedAt sql.NullTime
	var safeLayoutID sql.NullInt64
	var safeWaitingReason sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&res.ID, &res.Customer.ID, &res.Customer.Name, &res.Customer.Phone,
		&res.Customer.Address, &res.TotalQuantity, &res.ProductionType, &res.Status,
		&res.CreatedAt, &res.InternalSampleDeadline, &res.CustomerSampleDeadline,
		&safeActualFinished, &safeLayoutID, &safeWaitingReason, &safeInternalProductionDeadline,
		&safeCustomerProductionDeadline, &safeActualProductionStartedAt, &safeActualProductionCompletedAt,
	)
	if err != nil {
		return nil, errors.New("Order Not Found")
	}

	if safeActualFinished.Valid {
		res.ActualSampleFinishedAt = &safeActualFinished.Time
	}
	if safeInternalProductionDeadline.Valid {
		res.InternalProductionDeadline = &safeInternalProductionDeadline.Time
	}
	if safeCustomerProductionDeadline.Valid {
		res.CustomerProductionDeadline = &safeCustomerProductionDeadline.Time
	}
	if safeActualProductionStartedAt.Valid {
		res.ActualProductionStartedAt = &safeActualProductionStartedAt.Time
	}
	if safeActualProductionCompletedAt.Valid {
		res.ActualProductionCompletedAt = &safeActualProductionCompletedAt.Time
	}
	if safeWaitingReason.Valid {
		res.WaitingReason = &safeWaitingReason.String
	}
	if safeLayoutID.Valid {
		val := int(safeLayoutID.Int64)
		res.LayoutID = &val
	}

	//! fetching sizes
	res.Sizes = []OrderSize{}

	sizeQuery := `SELECT size_label, quantity FROM order_sizes WHERE order_id = $1`
	rows, err := r.db.Query(sizeQuery, id)
	if err != nil {
		fmt.Println("GetOrderByID size query error:", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var s OrderSize
			if err := rows.Scan(&s.SizeLabel, &s.Quantity); err != nil {
				fmt.Println("GetOrderByID Size scan error:", err)
			} else {
				res.Sizes = append(res.Sizes, s)
			}
		}

	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("Database connection dies during the loop")
	}

	//! fetch images
	imgQuery := `SELECT id, file_url FROM order_images WHERE order_id = $1`
	imgRows, err := r.db.Query(imgQuery, id)
	if err == nil {
		defer imgRows.Close()
		for imgRows.Next() {
			var img struct {
				ID       int    `json:"id"`
				ImageURL string `json:"image_url"`
			}
			if err := imgRows.Scan(&img.ID, &img.ImageURL); err != nil {
				res.Images = append(res.Images, img)
			}
		}
	}

	if err := imgRows.Err(); err != nil {
		return nil, errors.New("Database connection dies during the loop")
	}

	return &res, nil
}

// ! waiting for material
func (r *OrderRepository) SetWaitingForMaterials(orderID int, reason string) error {
	query := `UPDATE orders SET status = 'WAITING_FOR_MATERIALS', waiting_reason = $1 WHERE id = $2`
	_, err := r.db.Exec(query, reason, orderID)
	if err != nil {
		return errors.New("Failed to set waiting status and reason")
	}
	return nil
}

// ! resume order
func (r *OrderRepository) ResumeOrder(orderID int) error {
	query := `UPDATE orders SET status = 'READY_FOR_PRODUCTION', waiting_reason = NULL WHERE id = $1 AND status = 'WAITING_FOR_MATERIALS'`
	result, err := r.db.Exec(query, orderID)
	if err != nil {
		return errors.New("Failed to resume order")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("Order is not currently waiting for materials / not found")
	}

	return nil
}

// ! Hold Production
func (r *OrderRepository) HoldProduction(orderID int, reason string) error {
	query := `UPDATE orders SET waiting_reason = $1 WHERE id = $2`
	_, err := r.db.Exec(query, reason, orderID)
	return err
}

// ! Resume Production
func (r *OrderRepository) ResumeProduction(orderID int) error {
	query := `UPDATE orders SET waiting_reason = NULL WHERE id = $1`
	_, err := r.db.Exec(query, orderID)
	return err
}
