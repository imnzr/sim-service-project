package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/imnzr/sim-service-project/models"
)

type OrderRepository interface {
	CreateSimOrder(ctx context.Context, order *models.SimOrder) error
	GetSimOrderById(ctx context.Context, orderId uint) (*models.SimOrder, error)
	UpdateSimOrder(ctx context.Context, order *models.SimOrder) error
	GetSimOrderByInvoiceId(ctx context.Context, invoiceId string) (*models.SimOrder, error)
	GetUserSimOrders(ctx context.Context, userId uint, filter models.GetUserOrdersRequest) ([]models.SimOrder, error)
}

type OrderRepositoryImplement struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &OrderRepositoryImplement{
		Db: db,
	}
}

// CreateSimOrder implements OrderRepository.
func (o *OrderRepositoryImplement) CreateSimOrder(ctx context.Context, order *models.SimOrder) error {
	query := `
		INSERT INTO sim_orders(user_id,service, country, operator, price, invoice_id, sim_order_service_id, phone_number, otp, status, error_message, created_at, updated_at)
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?)
	`
	result, err := o.Db.ExecContext(ctx, query,
		order.UserId,
		order.Service,
		order.Country,
		order.Operator,
		order.Price,
		order.InvoiceId,
		order.SimOrderServiceId,
		order.PhoneNumber,
		order.OTP,
		order.Status,
		order.ErrorMessage,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		log.Println("failed to create order SIM")
		return fmt.Errorf("failed to create order SIM: %w", err)
	}
	id, err := result.LastInsertId()

	if err == nil {
		order.Id = uint(id)
	}

	return nil
}

// UpdateSimOrder implements OrderRepository : update status and data order in db
func (o *OrderRepositoryImplement) UpdateSimOrder(ctx context.Context, order *models.SimOrder) error {
	query := `
		UPDATE sim_orders SET
			user_id = ?, service = ?, country = ?, operator = ?, price = ?, invoice_id = ?,
			sim_order_service_id = ?, phone_number = ?, otp = ?, status = ?, error_message = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := o.Db.ExecContext(ctx, query,
		order.UserId,
		order.Service,
		order.Country,
		order.Operator,
		order.Price,
		order.InvoiceId,
		order.SimOrderServiceId,
		order.PhoneNumber,
		order.OTP,
		order.Status,
		order.ErrorMessage,
		order.UpdatedAt,
		order.Id,
	)
	if err != nil {
		log.Println("failed to update order SIM")
		return fmt.Errorf("failed to update order SIM: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		log.Println("order sim is not found for update")
		return fmt.Errorf("order sim is not found for update: %w", err)
	}

	return nil
}

// GetSimOrderById implements OrderRepository.
func (o *OrderRepositoryImplement) GetSimOrderById(ctx context.Context, orderId uint) (*models.SimOrder, error) {
	query := `
		SELECT id, user_id, service, country, operator, price, invoice_id,
			sim_order_service_id, phone_number, otp, status, error_message,
			created_at, updated_at
		FROM sim_orders WHERE id = ?
	`
	row := o.Db.QueryRowContext(ctx, query, orderId)

	var order models.SimOrder

	err := row.Scan(
		&order.Id,
		&order.UserId,
		&order.Service,
		&order.Country,
		&order.Operator,
		&order.Price,
		&order.InvoiceId,
		&order.SimOrderServiceId,
		&order.PhoneNumber,
		&order.OTP,
		&order.Status,
		&order.ErrorMessage,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order sim is not found")
		}
		return nil, fmt.Errorf("failed get order sim: %w", err)
	}
	return &order, nil
}

// GetSimOrderByInvoiceId implements OrderRepository.
func (o *OrderRepositoryImplement) GetSimOrderByInvoiceId(ctx context.Context, invoiceId string) (*models.SimOrder, error) {
	query := `
		SELECT id, user_id, service, country, operator, price, invoice_id, sim_order_service_id,
		phone_number, otp, status, error_message, created_at, updated_at
		FROM sim_orders WHERE invoice_id = ?
	`
	row := o.Db.QueryRowContext(ctx, query, invoiceId)

	var order models.SimOrder

	err := row.Scan(
		&order.Id,
		&order.UserId,
		&order.Service,
		&order.Country,
		&order.Operator,
		&order.Price,
		&order.InvoiceId,
		&order.SimOrderServiceId,
		&order.PhoneNumber,
		&order.OTP,
		&order.Status,
		&order.ErrorMessage,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order sim by invoice id is not found")
		}
		return nil, fmt.Errorf("failed get order sim by invoice id: %w", err)
	}

	return &order, nil
}

// GetUserSimOrders implements OrderRepository.
func (o *OrderRepositoryImplement) GetUserSimOrders(ctx context.Context, userId uint, filter models.GetUserOrdersRequest) ([]models.SimOrder, error) {
	baseQuery := `
		SELECT id, user_id, service, country, operator, price, invoice_id,
				sim_order_service_id, phone_number, otp, status, error_message,
				created_at, updated_at
		FROM sim_orders WHERE user_id = ?
	`

	args := []interface{}{userId}

	if filter.Status != "" {
		baseQuery += ` AND status = ?`
		args = append(args, filter.Status)
	}

	baseQuery += ` ORDER by created_at DESC`

	if filter.Limit > 0 {
		baseQuery += ` LIMIT ? `
		args = append(args, filter.Limit)
	}

	if filter.Offset >= 0 {
		baseQuery += ` OFFSET ?`
		args = append(args, filter.Offset)
	}

	rows, err := o.Db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed get list order SIM user: %w", err)
	}
	defer rows.Close()

	var orders []models.SimOrder

	for rows.Next() {
		var order models.SimOrder
		err := rows.Scan(
			&order.Id,
			&order.UserId,
			&order.Service,
			&order.Country,
			&order.Operator,
			&order.Price,
			&order.InvoiceId,
			&order.SimOrderServiceId,
			&order.PhoneNumber,
			&order.OTP,
			&order.Status,
			&order.ErrorMessage,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows order sim: %w", err)
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error from rows order sim: %w", err)
	}
	return orders, nil
}
