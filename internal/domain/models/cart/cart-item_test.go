package cart_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-go/internal/domain/models/cart"
	"github.com/stretchr/testify/assert"
)

func TestCartItem_NewCartItem_OnValidValues_ReturnsCartItem(t *testing.T) {
	productId := uuid.New()

	sut, err := cart.NewCartItem(productId, 2, 2550)

	assert.NoError(t, err)
	assert.Equal(t, productId, sut.ProductId)
	assert.Equal(t, int64(2550), sut.Price.Value)
	assert.Equal(t, int32(2), sut.Quantity.Value)
	assert.Equal(t, int64(5100), sut.TotalPrice().Value)
}

func TestCartItem_IncreaseQuantity_OnValidValues_UpdatesCartItem(t *testing.T) {
	productId := uuid.New()
	cart, _ := cart.NewCartItem(productId, 5, 2550)

	cart.IncreaseQuantity(12)
	cart.IncreaseQuantity(3)

	assert.Equal(t, productId, cart.ProductId)
	assert.Equal(t, int64(2550), cart.Price.Value)
	assert.Equal(t, int32(20), cart.Quantity.Value)
	assert.Equal(t, int64(51000), cart.TotalPrice().Value)
}

func TestCartItem_DecreaseQuantity_OnValidValues_UpdatesCartItem(t *testing.T) {
	productId := uuid.New()
	cart, _ := cart.NewCartItem(productId, 5, 2550)

	cart.IncreaseQuantity(12)
	cart.IncreaseQuantity(3)
	cart.IncreaseQuantity(5)
	cart.DecreaseQuantity(5)
	cart.IncreaseQuantity(8)
	cart.DecreaseQuantity(1)
	cart.DecreaseQuantity(2)

	assert.Equal(t, productId, cart.ProductId)
	assert.Equal(t, int64(2550), cart.Price.Value)
	assert.Equal(t, int32(25), cart.Quantity.Value)
	assert.Equal(t, int64(63750), cart.TotalPrice().Value)
}

func TestCartItem_DecreaseQuantity_OnDecreaseMoreThanCurrentQuantity_UpdatesQuantityToZero(t *testing.T) {
	productId := uuid.New()
	cart, _ := cart.NewCartItem(productId, 5, 2500)

	cart.DecreaseQuantity(10)

	assert.Equal(t, productId, cart.ProductId)
	assert.Equal(t, int64(0), cart.Price.Value)
	assert.Equal(t, int32(0), cart.Quantity.Value)
	assert.Equal(t, int64(0), cart.TotalPrice().Value)
}

func TestCartItem_NewCartItem_OnQuantityEqualsZero_ReturnsError(t *testing.T) {
	_, err := cart.NewCartItem(uuid.New(), 0, 2550)

	assert.EqualError(t, err, "cart item quantity cannot be less than one")
}

func TestCartItem_NewCartItem_OnNegativeQuantity_ReturnsError(t *testing.T) {
	_, err := cart.NewCartItem(uuid.New(), -1, 2550)

	assert.EqualError(t, err, "quantity value cannot be negative")
}

func TestCartItem_NewCartItem_OnNegativePrice_ReturnsError(t *testing.T) {
	_, err := cart.NewCartItem(uuid.New(), 1, -500)

	assert.EqualError(t, err, "money value cannot be negative")
}

func TestCartItem_IncreaseQuantity_OnQuantityEqualsZero_ReturnsError(t *testing.T) {
	cart, _ := cart.NewCartItem(uuid.New(), 0, 2500)

	err := cart.IncreaseQuantity(0)

	assert.EqualError(t, err, "cart item quantity cannot be less than one")
}

func TestCartItem_IncreaseQuantity_OnNegativeQuantity_ReturnsError(t *testing.T) {
	productId := uuid.New()
	cart, _ := cart.NewCartItem(productId, 5, 2550)

	err := cart.IncreaseQuantity(-4)

	assert.EqualError(t, err, "quantity value cannot be negative")
}

func TestCartItem_DecreaseQuantity_OnNegativeQuantity_ReturnsError(t *testing.T) {
	productId := uuid.New()
	cart, _ := cart.NewCartItem(productId, 5, 2550)

	err := cart.DecreaseQuantity(-4)

	assert.EqualError(t, err, "quantity value cannot be negative")
}

func TestCartItem_DecreaseQuantity_OnQuantityLessThanOne_ReturnsError(t *testing.T) {
	cart, _ := cart.NewCartItem(uuid.New(), 0, 2500)

	err := cart.DecreaseQuantity(0)

	assert.EqualError(t, err, "cart item quantity cannot be less than one")
}
