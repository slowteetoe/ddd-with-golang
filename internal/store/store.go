package store

import (
	"github.com/google/uuid"
	coffeeco "slowteetoe.com/coffeeco/internal"
)

type Store struct {
	ID              uuid.UUID
	Location        string
	ProductsForSale []coffeeco.Product
}
