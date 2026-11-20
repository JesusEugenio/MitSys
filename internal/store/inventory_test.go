package store

import (
	"context"
	"testing"

	"github.com/JesusEugenio/MitSys/internal/models"
	"github.com/shopspring/decimal"
)

func TestCreateProduct(t *testing.T) {
	testCases := []struct {
		name      string
		product   models.Product
		wantError bool
	}{
		{
			name: "success | create new product 'Arrachera'",
			product: models.Product{
				CutName:        "Arrachera",
				PricePerKilo:   decimal.NewFromFloat(231.24),
				AvailableStock: decimal.NewFromFloat(30.123),
			},
			wantError: false,
		},
		{
			name: "error | not correct number in PricePerKilo",
			product: models.Product{
				CutName:        "Solomillo",
				PricePerKilo:   decimal.NewFromFloat(2132123492083213.12),
				AvailableStock: decimal.NewFromFloat(21.213),
			},
			wantError: true,
		},
		{
			name: "error | not correct number in AvailableStock",
			product: models.Product{
				CutName:        "Milanesa de Pollo",
				PricePerKilo:   decimal.NewFromFloat(123.12),
				AvailableStock: decimal.NewFromFloat(1234349482734988230.123),
			},
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := testPool.Begin(context.Background())
			if err != nil {
				t.Fatalf("error initializing transaction for test: %v", err)
			}
			t.Cleanup(func() {
				err := tx.Rollback(context.Background())
				if err != nil {
					t.Logf("error making rollback in transaction: %v", err)
				}
			})

			store := NewInventory(tx)
			id, err := store.CreateProduct(context.Background(), tc.product)

			if (err != nil) != tc.wantError {
				t.Fatalf("CreateProduct() error = %v, wantErr = %v", err, tc.wantError)
			}

			if !tc.wantError && id == 0 {
				t.Fatal("CreateProduct() returned zero ID")
			}
		})
	}
}
