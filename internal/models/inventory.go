package models

import (
	"github.com/shopspring/decimal"
)

type Product struct {
	ID           int
	CutName      string
	PricePerKilo decimal.Decimal
	AvaibleStock decimal.Decimal
	IsActive     bool
}
