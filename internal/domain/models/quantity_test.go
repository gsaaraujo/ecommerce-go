package models_test

import (
	"testing"

	"github.com/gsaaraujo/ecommerce-go/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func Test_create_quantity_should_succeed(t *testing.T) {
	sut, err := models.NewQuantity(10)

	assert.NoError(t, err)
	assert.Equal(t, int16(10), sut.Value)
}

func Test_create_quantity_with_negative_value_should_fail(t *testing.T) {
	_, err := models.NewQuantity(-8)

	assert.EqualError(t, err, "quantity value cannot be negative")
}
