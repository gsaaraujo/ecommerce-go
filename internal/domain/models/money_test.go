package models_test

import (
	"testing"

	"github.com/gsaaraujo/ecommerce-go/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestMoney_NewMoney_OnValid_ReturnsMoney(t *testing.T) {
	sut, err := models.NewMoney(220)

	assert.NoError(t, err)
	assert.Equal(t, int64(220), sut.Value)
}

func TestMoney_NewMoney_OnNegativeValue_ReturnsError(t *testing.T) {
	_, err := models.NewMoney(-2)

	assert.EqualError(t, err, "money value cannot be negative")
}
