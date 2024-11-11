package models

import "errors"

type Quantity struct {
	Value int32
}

func NewQuantity(value int32) (Quantity, error) {
	if value < 0 {
		return Quantity{}, errors.New("quantity value cannot be negative")
	}

	return Quantity{
		Value: value,
	}, nil
}
