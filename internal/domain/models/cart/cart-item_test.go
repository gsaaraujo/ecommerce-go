package cart_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-go/internal/domain/models/cart"
	"github.com/stretchr/testify/assert"
)

func Test_create_cart_item_should_succeed(t *testing.T) {
	productId := uuid.New()

	sut, err := cart.NewCartItem(productId, 2, 2550)

	assert.NoError(t, err)
	assert.Equal(t, productId, sut.ProductId)
	assert.Equal(t, int64(2550), sut.Price.Value)
	assert.Equal(t, int16(2), sut.Quantity.Value)
	assert.Equal(t, int64(5100), sut.TotalPrice().Value)
}

func Test_increase_quantity_should_succeed(t *testing.T) {
	productId := uuid.New()
	cart, _ := cart.NewCartItem(productId, 5, 2550)

	cart.IncreaseQuantity(12)
	cart.IncreaseQuantity(3)

	assert.Equal(t, productId, cart.ProductId)
	assert.Equal(t, int64(2550), cart.Price.Value)
	assert.Equal(t, int16(20), cart.Quantity.Value)
	assert.Equal(t, int64(51000), cart.TotalPrice().Value)
}

func Test_increase_and_decrease_quantity_should_succeed(t *testing.T) {
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
	assert.Equal(t, int16(25), cart.Quantity.Value)
	assert.Equal(t, int64(63750), cart.TotalPrice().Value)
}

func Test_decrease_quantity_with_higher_value_than_current_should_result_in_price_and_quantity_equals_zero(t *testing.T) {
	productId := uuid.New()
	cart, _ := cart.NewCartItem(productId, 5, 2500)

	cart.DecreaseQuantity(10)

	assert.Equal(t, productId, cart.ProductId)
	assert.Equal(t, int64(0), cart.Price.Value)
	assert.Equal(t, int16(0), cart.Quantity.Value)
	assert.Equal(t, int64(0), cart.TotalPrice().Value)
}

func Test_create_cart_item_with_quantity_equals_zero_should_fail(t *testing.T) {
	_, err := cart.NewCartItem(uuid.New(), 0, 2550)

	assert.EqualError(t, err, "cart item quantity cannot be less than one")
}

func Test_create_cart_item_with_negative_quantity_should_fail(t *testing.T) {
	_, err := cart.NewCartItem(uuid.New(), -1, 2550)

	assert.EqualError(t, err, "quantity value cannot be negative")
}

func Test_create_cart_item_with_negative_price_should_fail(t *testing.T) {
	_, err := cart.NewCartItem(uuid.New(), 1, -500)

	assert.EqualError(t, err, "money value cannot be negative")
}

func Test_increase_quantity_with_quantity_less_than_one_should_fail(t *testing.T) {
	cart, _ := cart.NewCartItem(uuid.New(), 0, 2500)

	err := cart.IncreaseQuantity(0)

	assert.EqualError(t, err, "cart item quantity cannot be less than one")
}

func Test_increase_quantity_with_negative_quantity_should_fail(t *testing.T) {
	productId := uuid.New()
	cart, _ := cart.NewCartItem(productId, 5, 2550)

	err := cart.IncreaseQuantity(-4)

	assert.EqualError(t, err, "quantity value cannot be negative")
}

func Test_decrease_quantity_with_negative_quantity_should_fail(t *testing.T) {
	productId := uuid.New()
	cart, _ := cart.NewCartItem(productId, 5, 2550)

	err := cart.DecreaseQuantity(-4)

	assert.EqualError(t, err, "quantity value cannot be negative")
}

func Test_decrease_quantity_with_quantity_less_than_one_should_fail(t *testing.T) {
	cart, _ := cart.NewCartItem(uuid.New(), 0, 2500)

	err := cart.DecreaseQuantity(0)

	assert.EqualError(t, err, "cart item quantity cannot be less than one")
}
