package models

import "errors"

type Money struct {
	Value int64
}

func NewMoney(value int64) (Money, error) {
	if value < 0 {
		return Money{}, errors.New("money value cannot be negative")
	}

	return Money{
		Value: value,
	}, nil
}
