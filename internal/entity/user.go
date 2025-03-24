package entity

import (
	"github.com/gofrs/uuid/v5"
	"github.com/shopspring/decimal"
)

type User struct {
	ID      uuid.UUID       `json:"id"`
	Name    string          `json:"name"`
	Email   string          `json:"email"`
	Age     int             `json:"age"`
	Balance decimal.Decimal `json:"balance"`
}
