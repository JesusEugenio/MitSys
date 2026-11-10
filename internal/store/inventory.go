package store

import (
	"context"
	"fmt"
	"strconv"

	"github.com/JesusEugenio/MitSys/internal/models"
	"github.com/shopspring/decimal"
)

type InvetoryStore interface {
	CreateProduct(ctx context.Context, p models.Product) (int, error)
	UpdateProductPrice(ctx context.Context, id int, newPrice decimal.Decimal) error
	SearchProducts(ctx context.Context, search string) ([]models.Product, error)
}

type inventoryStore struct {
	db DBTX
}

func NewInventory(db DBTX) InvetoryStore {
	return &inventoryStore{
		db: db,
	}
}

func (s *inventoryStore) CreateProduct(ctx context.Context, p models.Product) (int, error) {
	var id int

	query := `INSERT INTO products (cut_name, price_per_kilo, available_stock) VALUES ($1, $2, $3) RETURNING id`

	err := s.db.QueryRow(ctx, query, p.CutName, p.PricePerKilo, p.AvaibleStock).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error creating a new product: %w", err)
	}

	return id, nil
}

func (s *inventoryStore) UpdateProductPrice(ctx context.Context, id int, newPrice decimal.Decimal) error {
	query := `UPDATE products SET price_per_kilo = $1 WHERE id = $2 AND is_active = TRUE`

	result, err := s.db.Exec(ctx, query, newPrice, id)
	if err != nil {
		return fmt.Errorf("error updating product price: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (s *inventoryStore) SearchProducts(ctx context.Context, search string) ([]models.Product, error) {
	query := `SELECT id, cut_name, price_per_kilo, available_stock, is_active FROM products WHERE is_active = TRUE`
	var args []any

	if search != "" {
		if id, err := strconv.Atoi(search); err == nil {
			query += ` AND (cut_name ILIKE $1 OR id = $2)`
			args = append(args, "%"+search+"%", id)
		} else {
			query += ` AND cut_name ILIKE $1`
			args = append(args, "%"+search+"%")
		}
	}

	query += ` ORDER BY id ASC`

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.CutName, &p.PricePerKilo, &p.AvaibleStock, &p.IsActive)
		if err != nil {
			return nil, fmt.Errorf("error scanning product row: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product rows: %w", err)
	}

	return products, nil
}
