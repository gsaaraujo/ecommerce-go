package usecases

import (
	"errors"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-go/internal/application/gateways"
	"github.com/gsaaraujo/ecommerce-go/internal/application/repositories"
	"github.com/gsaaraujo/ecommerce-go/internal/domain/models/cart"
)

type AddProductToCartInput struct {
	CustomerId uuid.UUID
	ProductId  uuid.UUID
	Price      int64
	Quantity   int16
}

type IAddProductToCart interface {
	Execute(input AddProductToCartInput) error
}

type AddProductToCart struct {
	CustomerGateway gateways.ICustomerGateway
	CartRepository  repositories.ICartRepository
}

func (a *AddProductToCart) Execute(input AddProductToCartInput) error {
	customerExists, err := a.CustomerGateway.ExistsByCustomerId(input.CustomerId)
	if err != nil {
		return err
	}

	if !customerExists {
		return errors.New("customer not found")
	}

	customerCart, err := a.CartRepository.FindOneByCustomerId(input.CustomerId)
	if err != nil {
		return err
	}

	if customerCart != nil {
		err := customerCart.AddItem(input.ProductId, input.Quantity, input.Price)
		if err != nil {
			return err
		}

		err = a.CartRepository.Update(*customerCart)
		if err != nil {
			return err
		}

		return nil
	}

	newCart, err := cart.NewCart(input.CustomerId)
	if err != nil {
		return err
	}

	err = newCart.AddItem(input.ProductId, input.Quantity, input.Price)
	if err != nil {
		return err
	}

	err = a.CartRepository.Create(newCart)
	if err != nil {
		return err
	}

	return nil
}
