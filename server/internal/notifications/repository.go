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
	query := `SELECT id, status, internal_production_deadline
						FROM orders WHERE status NOT IN ('SHIPPED', 'READY_FOR_SHIPPING')
						AND internal_production_deadline < CURRENT_TIMESTAMP
						AND id NOT IN (SELECT order_id FROM notifications WHERE alert_type = 'PRODUCTION_DELAY)`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var delayed []DelayedOrder
	for rows.Next() {
		var order DelayedOrder
		if err := rows.Scan(&order.OrderID, &order.Status, &order.Deadline); err == nil {
			delayed = append(delayed, order)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("Database Connection dies during the loop")
	}

	return delayed, nil
}

func (r *NotificationRepository) MarkAlertSent(orderID int, AlertType string) (int, error) {
	query := `INSERT INTO notifications (order_id, alert_type) VALUES ($1, $2) RETURNING id`

	var newID int
	err := r.db.QueryRow(query, orderID, AlertType).Scan(&newID)
	return newID, err
}
