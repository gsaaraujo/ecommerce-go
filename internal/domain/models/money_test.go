package models_test

import (
	"testing"

	"github.com/gsaaraujo/ecommerce-go/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func Test_create_money_should_succeed(t *testing.T) {
	sut, err := models.NewMoney(220)

	assert.NoError(t, err)
	assert.Equal(t, int64(220), sut.Value)
}

func Test_create_money_with_negative_value_should_fail(t *testing.T) {
	_, err := models.NewMoney(-2)

	assert.EqualError(t, err, "money value cannot be negative")
}
