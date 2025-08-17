package repository

import (
	"context"
	"database/sql"

	"github.com/imnzr/sim-service-project/models"
)

type SimOrderRepository interface {
	Create(ctx context.Context, order *models.SimOrder) (int, error)
	UpdateStatusByInvoiceId(ctx context.Context, invoiceId, status string) error
	UpdateAfterPayment(ctx context.Context, invoiceId string, simServiceId int, phoneNumber string) error
	GetByInvoiceId(ctx context.Context, invoiceId string) (*models.SimOrder, error)
	GetById(ctx context.Context, id int) (*models.SimOrder, error)
	AttachSimDataService(ctx context.Context, orderId int, data *models.ResponsOrderFromService) error
}

type SimOrderImplement struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) SimOrderRepository {
	return &SimOrderImplement{
		db: db,
	}
}

// GetById implements SimOrderRepository.
func (s *SimOrderImplement) GetById(ctx context.Context, id int) (*models.SimOrder, error) {
	query := "SELECT * FROM sim_orders WHERE id = ?"
	row := s.db.QueryRowContext(ctx, query, id)

	var order models.SimOrder

	err := row.Scan(
		&order.Id,
		&order.UserId,
		&order.Service,
		&order.Country,
		&order.Operator,
		&order.PriceSell,
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
		if err == sql.ErrNoRows {
			return nil, nil // No order found
		}
		return nil, err // Other error
	}
	return &order, nil // Return the found order
}

// AttachSimDataService implements SimOrderRepository.
func (s *SimOrderImplement) AttachSimDataService(ctx context.Context, orderId int, data *models.ResponsOrderFromService) error {
	query := `
		UPDATE sim_orders
		SET sim_order_service_id = ?, phone_number = ?, status = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := s.db.ExecContext(ctx, query,
		data.Id,
		data.Phone,
		data.Status,
		orderId,
	)
	return err
}

// CreateOrder implements SimOrderRepository.
func (s *SimOrderImplement) Create(ctx context.Context, order *models.SimOrder) (int, error) {
	query := `
		INSERT INTO sim_orders(user_id, service, country, operator, price, invoice_id, status)
		VALUES(?,?,?,?,?,?,?)
	`
	result, err := s.db.ExecContext(ctx, query,
		order.UserId,
		order.Service,
		order.Country,
		order.Operator,
		order.PriceSell,
		order.InvoiceId,
		order.Status,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

// GetByInvoiceId implements SimOrderRepository.
func (s *SimOrderImplement) GetByInvoiceId(ctx context.Context, invoiceId string) (*models.SimOrder, error) {
	query := `
		SELECT * FROM sim_orders WHERE invoice_id = ?
	`
	row := s.db.QueryRowContext(ctx, query, invoiceId)

	var order models.SimOrder
	err := row.Scan(
		&order.Id,
		&order.UserId,
		&order.Service,
		&order.Country,
		&order.Operator,
		&order.PriceSell,
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
		return nil, err
	}
	return &order, nil
}

// UpdateAfterPayment implements SimOrderRepository.
func (s *SimOrderImplement) UpdateAfterPayment(ctx context.Context, invoiceId string, simServiceId int, phoneNumber string) error {
	query := `
		UPDATE sim_orders SET sim_order_service_id = ?, phone_number = ?, status = 'ACTIVE'
		WHERE invoice_id = ?
	`
	_, err := s.db.ExecContext(ctx, query, simServiceId, phoneNumber, invoiceId)
	return err
}

// UpdateStatusByInvoiceId implements SimOrderRepository.
func (s *SimOrderImplement) UpdateStatusByInvoiceId(ctx context.Context, invoiceId string, status string) error {
	query := `
		UPDATE sim_orders SET status = ? WHERE invoice_id = ?
	`
	_, err := s.db.ExecContext(ctx, query, status, invoiceId)

	return err
}
