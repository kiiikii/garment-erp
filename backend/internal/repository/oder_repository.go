package repository

import (
	"database/sql"

	"github.com/kiiikii/garment-erp/backend/internal/models"
)

type OrderRepo struct {
	DB *sql.DB
}

func (r *OrderRepo) InsertOrder(req models.CreateOrderReq) (int, error) {
	var newID int
	query := `INSERT INTO orders (customer_name, product_type) VALUES ($1, $2) RETURNING id`
	err := r.DB.QueryRow(query, req.CustName, req.ProdType).Scan(&newID)
	return newID, err
}

func (r *OrderRepo) GetAllOrders() ([]models.Order, error) {
	query := `SELECT id, customer_name, product_type, status FROM orders ORDER BY id DESC`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := []models.Order{}

	for rows.Next() {
		var o models.Order

		err := rows.Scan(&o.ID, &o.CustName, &o.ProdType, &o.Status)
		if err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepo) UpdateOrderStatus(id int, newStatus string) error {
	_, err := r.DB.Exec(`UPDATE orders SET status = $1 WHERE id = $2`, newStatus, id)
	return err
}
