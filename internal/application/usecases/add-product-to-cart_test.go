package usecases_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-go/internal/application/gateways"
	"github.com/gsaaraujo/ecommerce-go/internal/application/usecases"
	"github.com/gsaaraujo/ecommerce-go/internal/domain/models/cart"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CustomerGatewayMock struct {
	mock.Mock
}

func (c *CustomerGatewayMock) ExistsById(customerId uuid.UUID) (bool, error) {
	args := c.Called(customerId)
	return args.Bool(0), args.Error(1)
}

type CartRepositoryMock struct {
	mock.Mock
}

func (c *CartRepositoryMock) Create(cart cart.Cart) error {
	args := c.Called(cart)
	return args.Error(0)
}

func (c *CartRepositoryMock) Update(cart cart.Cart) error {
	args := c.Called(cart)
	return args.Error(0)
}

func (c *CartRepositoryMock) FindOneByCustomerId(customerId uuid.UUID) (*cart.Cart, error) {
	args := c.Called(customerId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cart.Cart), args.Error(1)
}

type ProductGatewayMock struct {
	mock.Mock
}

func (p *ProductGatewayMock) FindOneById(id uuid.UUID) (*gateways.ProductDTO, error) {
	args := p.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gateways.ProductDTO), args.Error(1)
}

type AddProductToCartSuite struct {
	suite.Suite
	addProductToCart    usecases.AddProductToCart
	customerGatewayMock CustomerGatewayMock
	productGatewayMock  ProductGatewayMock
	cartRepositoryMock  CartRepositoryMock
}

func (a *AddProductToCartSuite) SetupTest() {
	a.customerGatewayMock = CustomerGatewayMock{}
	a.productGatewayMock = ProductGatewayMock{}
	a.cartRepositoryMock = CartRepositoryMock{}

	a.addProductToCart = usecases.AddProductToCart{
		CustomerGateway: &a.customerGatewayMock,
		ProductGateway:  &a.productGatewayMock,
		CartRepository:  &a.cartRepositoryMock,
	}
}

func (a *AddProductToCartSuite) Test_add_product_to_new_cart_should_succeed() {
	product := gateways.ProductDTO{
		Id:    uuid.New(),
		Price: int64(3550),
	}
	a.customerGatewayMock.On("ExistsById", mock.Anything).Return(true, nil)
	a.cartRepositoryMock.On("FindOneByCustomerId", mock.Anything).Return(nil, nil)
	a.productGatewayMock.On("FindOneById", mock.Anything).Return(&product, nil)
	a.cartRepositoryMock.On("Create", mock.Anything).Return(nil)
	input := usecases.AddProductToCartInput{
		CustomerId: uuid.New(),
		ProductId:  uuid.New(),
		Quantity:   int32(3),
	}

	err := a.addProductToCart.Execute(input)

	a.Equal(nil, err)
	a.cartRepositoryMock.AssertNumberOfCalls(a.T(), "Create", 1)
	a.cartRepositoryMock.AssertNumberOfCalls(a.T(), "Update", 0)
}

func (a *AddProductToCartSuite) Test_add_product_to_existing_cart_should_succeed() {
	customerCart := cart.Cart{
		Id:         uuid.New(),
		CustomerId: uuid.New(),
		Items:      []cart.CartItem{},
	}
	product := gateways.ProductDTO{
		Id:    uuid.New(),
		Price: int64(3550),
	}
	a.customerGatewayMock.On("ExistsById", mock.Anything).Return(true, nil)
	a.cartRepositoryMock.On("FindOneByCustomerId", mock.Anything).Return(&customerCart, nil)
	a.productGatewayMock.On("FindOneById", mock.Anything).Return(&product, nil)
	a.cartRepositoryMock.On("Update", mock.Anything).Return(nil)
	input := usecases.AddProductToCartInput{
		CustomerId: uuid.New(),
		ProductId:  uuid.New(),
		Quantity:   int32(3),
	}

	err := a.addProductToCart.Execute(input)

	a.Equal(nil, err)
	a.cartRepositoryMock.AssertNumberOfCalls(a.T(), "Create", 0)
	a.cartRepositoryMock.AssertNumberOfCalls(a.T(), "Update", 1)
}

func (a *AddProductToCartSuite) Test_add_product_to_cart_with_customer_that_does_not_exist_should_fail() {
	product := gateways.ProductDTO{
		Id:    uuid.New(),
		Price: int64(3550),
	}
	a.productGatewayMock.On("FindOneById", mock.Anything).Return(&product, nil)
	a.customerGatewayMock.On("ExistsById", mock.Anything).Return(false, nil)
	input := usecases.AddProductToCartInput{
		CustomerId: uuid.New(),
		ProductId:  uuid.New(),
		Quantity:   int32(3),
	}

	err := a.addProductToCart.Execute(input)

	a.EqualError(err, "customer not found")
}

func (a *AddProductToCartSuite) Test_add_product_to_cart_with_product_that_does_not_exist_should_fail() {
	a.customerGatewayMock.On("ExistsById", mock.Anything).Return(true, nil)
	a.productGatewayMock.On("FindOneById", mock.Anything).Return(nil, nil)
	input := usecases.AddProductToCartInput{
		CustomerId: uuid.New(),
		ProductId:  uuid.New(),
		Quantity:   int32(3),
	}

	err := a.addProductToCart.Execute(input)

	a.EqualError(err, "product not found")
}

func TestAddProductToCart(t *testing.T) {
	suite.Run(t, new(AddProductToCartSuite))
}
