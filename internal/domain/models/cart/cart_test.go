package cart_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-go/internal/domain/models/cart"
	"github.com/stretchr/testify/assert"
)

func TestCart_NewCart_OnValidValues_ReturnsCart(t *testing.T) {
	customerId := uuid.New()

	sut, err := cart.NewCart(customerId)

	assert.NoError(t, err)
	assert.Equal(t, customerId, sut.CustomerId)
	assert.Equal(t, []cart.CartItem{}, sut.Items)
}

func TestCart_AddItem_OnAddingSameProduct_UpdatesCart(t *testing.T) {
	customerId := uuid.New()
	product1 := uuid.New()
	cart, _ := cart.NewCart(customerId)

	cart.AddItem(product1, 2, 32000)
	cart.AddItem(product1, 5, 32000)

	assert.Equal(t, int(1), len(cart.Items))
	assert.Equal(t, int32(7), cart.TotalQuantity().Value)
	assert.Equal(t, int64(224000), cart.TotalPrice().Value)
}

func TestCart_AddItem_OnAddingDifferentProducts_UpdatesCart(t *testing.T) {
	customerId := uuid.New()
	product1 := uuid.New()
	product2 := uuid.New()
	product3 := uuid.New()
	cart, _ := cart.NewCart(customerId)

	cart.AddItem(product1, 2, 32000)
	cart.AddItem(product2, 5, 17340)
	cart.AddItem(product2, 1, 17340)
	cart.AddItem(product2, 4, 17340)
	cart.AddItem(product3, 9, 1550)

	assert.Equal(t, int(3), len(cart.Items))
	assert.Equal(t, int32(21), cart.TotalQuantity().Value)
	assert.Equal(t, int64(251350), cart.TotalPrice().Value)
}

func TestCart_RemoveItem_OnAddingAndRemovingSameProduct_UpdatesCart(t *testing.T) {
	customerId := uuid.New()
	product1 := uuid.New()
	cart, _ := cart.NewCart(customerId)

	cart.AddItem(product1, 2, 32000)
	cart.RemoveItem(product1)

	assert.Equal(t, int(0), len(cart.Items))
	assert.Equal(t, int32(0), cart.TotalQuantity().Value)
	assert.Equal(t, int64(0), cart.TotalPrice().Value)
}

func TestCart_RemoveItem_OnAddingAndRemovingDifferentProducts_UpdatesCart(t *testing.T) {
	customerId := uuid.New()
	product1 := uuid.New()
	product2 := uuid.New()
	product3 := uuid.New()
	cart, _ := cart.NewCart(customerId)

	cart.AddItem(product1, 2, 32000)
	cart.AddItem(product2, 5, 17340)
	cart.AddItem(product2, 1, 17340)
	cart.AddItem(product2, 4, 17340)
	cart.AddItem(product3, 9, 1550)
	cart.RemoveItem(product2)

	assert.Equal(t, int(2), len(cart.Items))
	assert.Equal(t, int32(11), cart.TotalQuantity().Value)
	assert.Equal(t, int64(77950), cart.TotalPrice().Value)
}

func TestCart_AddItem_OnNegativeQuantity_ReturnsError(t *testing.T) {
	customerId := uuid.New()
	cart, _ := cart.NewCart(customerId)

	err := cart.AddItem(uuid.New(), -2, 32000)

	assert.EqualError(t, err, "quantity value cannot be negative")
}

func TestCart_AddItem_OnQuantityEqualsZero_ReturnsError(t *testing.T) {
	customerId := uuid.New()
	cart, _ := cart.NewCart(customerId)

	err := cart.AddItem(uuid.New(), 0, 32000)

	assert.EqualError(t, err, "cart item quantity cannot be less than one")
}

func TestCart_AddItem_OnNegativePrice_ReturnsError(t *testing.T) {
	customerId := uuid.New()
	cart, _ := cart.NewCart(customerId)

	err := cart.AddItem(uuid.New(), 2, -550)

	assert.EqualError(t, err, "money value cannot be negative")
}

func TestCart_RemoveItem_OnProdutNotInCart_ReturnsError(t *testing.T) {
	customerId := uuid.New()
	product1 := uuid.New()
	product2 := uuid.New()
	cart, _ := cart.NewCart(customerId)

	cart.AddItem(product1, 2, 32000)
	err := cart.RemoveItem(product2)

	assert.EqualError(t, err, "product not found in cart")
}

func TestCart_RemoveItem_OnCartEmpty_ReturnsError(t *testing.T) {
	cart, _ := cart.NewCart(uuid.New())

	err := cart.RemoveItem(uuid.New())

	assert.EqualError(t, err, "cart is empty")
}
