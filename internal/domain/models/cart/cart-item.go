package cart

import (
	"errors"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-go/internal/domain/models"
)

type CartItem struct {
	Id        uuid.UUID
	ProductId uuid.UUID
	Quantity  models.Quantity
	Price     models.Money
}

func NewCartItem(productId uuid.UUID, quantity int16, price int64) (CartItem, error) {
	if _, err := models.NewQuantity(quantity); err != nil {
		return CartItem{}, err
	}

	if _, err := models.NewMoney(price); err != nil {
		return CartItem{}, err
	}

	MINIMUM_QUANTITY := int16(1)
	if quantity < MINIMUM_QUANTITY {
		return CartItem{}, errors.New("cart item quantity cannot be less than one")
	}

	return CartItem{
		Id:        uuid.New(),
		ProductId: productId,
		Quantity:  models.Quantity{Value: quantity},
		Price:     models.Money{Value: price},
	}, nil
}

func (c *CartItem) IncreaseQuantity(quantity int16) error {
	if _, err := models.NewQuantity(quantity); err != nil {
		return err
	}

	MINIMUM_QUANTITY := int16(1)
	if quantity < MINIMUM_QUANTITY {
		return errors.New("cart item quantity cannot be less than one")
	}

	c.Quantity = models.Quantity{Value: c.Quantity.Value + quantity}
	return nil
}

func (c *CartItem) DecreaseQuantity(quantity int16) error {
	if _, err := models.NewQuantity(quantity); err != nil {
		return err
	}

	MINIMUM_QUANTITY := int16(1)
	if quantity < MINIMUM_QUANTITY {
		return errors.New("cart item quantity cannot be less than one")
	}

	difference := c.Quantity.Value - quantity
	if difference < 0 {
		c.Quantity = models.Quantity{Value: 0}
		c.Price = models.Money{Value: 0}
		return nil
	}

	c.Quantity = models.Quantity{Value: difference}
	return nil
}

func (c *CartItem) TotalPrice() models.Money {
	return models.Money{
		Value: c.Price.Value * int64(c.Quantity.Value),
	}
}
