package service

import (
	"context"
	"fmt"

	"github.com/JesusEugenio/MitSys/internal/models"
	"github.com/JesusEugenio/MitSys/internal/store"
	"github.com/shopspring/decimal"
)

type InventoryService interface {
	CreateProduct(ctx context.Context, p models.Product) (int, error)
	UpdateProductPrice(ctx context.Context, id int, newPrice decimal.Decimal) error
	SearchProducts(ctx context.Context, search string) ([]models.Product, error)
	ViewProduct(ctx context.Context, id int) (models.Product, error)
	UpdateProduct(ctx context.Context, id int, p models.Product) error
}

type inventoryService struct {
	InventoryStore store.InvetoryStore
}

func NewInvetoryService(s store.InvetoryStore) InventoryService {
	return &inventoryService{InventoryStore: s}
}

func (s inventoryService) CreateProduct(ctx context.Context, p models.Product) (int, error) {

	//
	// Validations ..
	//

	id, err := s.InventoryStore.CreateProduct(ctx, p)
	if err != nil {
		return -1, fmt.Errorf("error creating product in datastore: %w", err)
	}

	return id, nil
}

func (s inventoryService) UpdateProductPrice(ctx context.Context, id int, newPrice decimal.Decimal) error {

	err := s.InventoryStore.UpdateProductPrice(ctx, id, newPrice)
	if err != nil {
		return fmt.Errorf("error updating price in datastore: %w", err)
	}

	return nil
}

func (s inventoryService) SearchProducts(ctx context.Context, search string) ([]models.Product, error) {
	products, err := s.InventoryStore.SearchProducts(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("error searching products in datastore: %w", err)
	}

	return products, nil
}

func (s inventoryService) ViewProduct(ctx context.Context, id int) (models.Product, error) {
	product, err := s.InventoryStore.ViewProduct(ctx, id)
	if err != nil {
		return models.Product{}, fmt.Errorf("error searching for the product: %w", err)
	}

	return product, nil
}

func (s inventoryService) UpdateProduct(ctx context.Context, id int, p models.Product) error {
	err := s.InventoryStore.UpdateProduct(ctx, id, p)

	if err != nil {
		return fmt.Errorf("error updating product in datastore: %w", err)
	}

	return nil
}
