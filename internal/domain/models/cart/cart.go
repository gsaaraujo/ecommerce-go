package cart

import (
	"errors"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-go/internal/domain/models"
)

type Cart struct {
	Id         uuid.UUID
	CustomerId uuid.UUID
	Items      []CartItem
}

func NewCart(customerId uuid.UUID) (Cart, error) {
	return Cart{
		Id:         uuid.New(),
		CustomerId: customerId,
		Items:      []CartItem{},
	}, nil
}

func (c *Cart) AddItem(productId uuid.UUID, quantity int32, price int64) error {
	if _, err := models.NewQuantity(quantity); err != nil {
		return err
	}

	if _, err := models.NewMoney(price); err != nil {
		return err
	}

	for i, item := range c.Items {
		if item.ProductId == productId {
			c.Items[i].IncreaseQuantity(quantity)
			return nil
		}
	}

	cartItem, err := NewCartItem(productId, quantity, price)

	if err != nil {
		return err
	}

	c.Items = append(c.Items, cartItem)
	return nil
}

func (c *Cart) RemoveItem(productId uuid.UUID) error {
	if len(c.Items) == 0 {
		return errors.New("cart is empty")
	}

	for i, item := range c.Items {
		if item.ProductId == productId {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return nil
		}
	}

	return errors.New("product not found in cart")
}

func (c *Cart) TotalQuantity() models.Quantity {
	totalQuantity := int32(0)

	for _, item := range c.Items {
		totalQuantity += item.Quantity.Value
	}

	return models.Quantity{
		Value: totalQuantity,
	}
}

func (c *Cart) TotalPrice() models.Money {
	totalPrice := int64(0)

	for _, item := range c.Items {
		totalPrice += item.TotalPrice().Value
	}

	return models.Money{
		Value: totalPrice,
	}
}
