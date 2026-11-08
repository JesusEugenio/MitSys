package store

import (
	"context"
	"fmt"

	"github.com/JesusEugenio/MitSys/internal/models"
)

type InvetoryStore interface {
	CreateProduct(ctx context.Context, p models.Product) (int, error)
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
