package notifications

import (
	"database/sql"
	"errors"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) GetDelayedProductionOrder() ([]DelayedOrder, error) {
	query := `SELECT o.id, o.status, o.internal_production_deadline, c.email, c.phone FROM orders o
						JOIN customers c ON o.customer_id = c.id WHERE o.status NOT IN ('SHIPPED', 'READY_FOR_SHIPPING')
						AND o.internal_production_deadline < CURRENT_TIMESTAMP AND o.id NOT IN (SELECT order_id FROM notifications
						WHERE alert_type = 'PRODUCTION_DELAYED)`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var delayed []DelayedOrder
	for rows.Next() {
		var order DelayedOrder
		if err := rows.Scan(&order.OrderID, &order.Status, &order.Deadline, &order.CustomerEmail, &order.CustomerPhone); err == nil {
			delayed = append(delayed, order)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("Database connection dies during the loop")
	}
	return delayed, nil
}

func (r *NotificationRepository) MarkAlertSent(orderID int, AlertType string) (int, error) {
	query := `INSERT INTO notifications (order_id, alert_type) VALUES ($1, $2) RETURNING id`

	var newID int
	err := r.db.QueryRow(query, orderID, AlertType).Scan(&newID)
	return newID, err
}
