package models_test

import (
	"testing"

	"github.com/gsaaraujo/ecommerce-go/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestQuantity_NewQuantity_OnValidValue_ReturnsQuantity(t *testing.T) {
	sut, err := models.NewQuantity(10)

	assert.NoError(t, err)
	assert.Equal(t, int32(10), sut.Value)
}

func TestQuantity_NewQuantity_OnNegativeValue_ReturnsError(t *testing.T) {
	_, err := models.NewQuantity(-8)

	assert.EqualError(t, err, "quantity value cannot be negative")
}
