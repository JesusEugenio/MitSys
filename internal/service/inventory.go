package service

import (
	"context"
	"fmt"

	"github.com/JesusEugenio/MitSys/internal/models"
	"github.com/JesusEugenio/MitSys/internal/store"
)

type InventoryService interface {
	CreateProduct(ctx context.Context, p models.Product) (int, error)
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
