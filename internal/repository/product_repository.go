package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/imnzr/sim-service-project/models"
)

type ProductRepository interface {
	Upsert(ctx context.Context, product *models.SimProduct) error
	FindByKey(ctx context.Context, service, country, operator string) (*models.SimProduct, error)
}

type ProductImplementation struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &ProductImplementation{
		db: db,
	}
}

// Upsert implements ProductRepository.
func (p *ProductImplementation) Upsert(ctx context.Context, product *models.SimProduct) error {
	query := "INSERT INTO product(service, country, operator, price_default, price_sell) VALUES(?,?,?,?,?)"

	result, err := p.db.ExecContext(ctx, query,
		product.Service,
		product.Country,
		product.Operator,
		product.PriceDefault,
		product.PriceSell,
	)

	if err != nil {
		log.Printf("failed to insert product: %v\n", err)
		return err
	}

	rows, _ := result.RowsAffected()

	fmt.Printf("Upsert success: %s | %s | %s (affected: %d)\n",
		product.Service,
		product.Country,
		product.Operator,
		rows,
	)

	return err
}

// FindByKey implements ProductRepository.
func (p *ProductImplementation) FindByKey(ctx context.Context, service string, country string, operator string) (*models.SimProduct, error) {
	query := `
		SELECT service, country, operator, price_sell
		FROM product
		WHERE service = ? AND country = ? AND operator = ?
	`
	var product models.SimProduct

	err := p.db.QueryRowContext(ctx, query, service, country, operator).Scan(
		&product.Service,
		&product.Country,
		&product.Operator,
		&product.PriceSell,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}
