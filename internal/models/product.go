package models

import "github.com/shopspring/decimal"

type Product struct {
	ID            int
	CutName       string
	PricePerKilo  decimal.Decimal
	Avaible_stock decimal.Decimal
	IsActive      bool
}
