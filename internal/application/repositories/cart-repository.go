package repositories

import (
	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-go/internal/domain/models/cart"
)

type ICartRepository interface {
	Create(cart cart.Cart) error
	Update(cart cart.Cart) error
	FindOneByCustomerId(customerId uuid.UUID) (*cart.Cart, error)
}
