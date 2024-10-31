package usecases_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-go/internal/application/usecases"
	"github.com/gsaaraujo/ecommerce-go/internal/domain/models/cart"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CustomerGatewayMock struct {
	mock.Mock
}

type CartRepositoryMock struct {
	mock.Mock
}

func (m *CustomerGatewayMock) ExistsByCustomerId(customerId uuid.UUID) (bool, error) {
	args := m.Called(customerId)
	return args.Bool(0), args.Error(1)
}

func (m *CartRepositoryMock) Create(cart cart.Cart) error {
	args := m.Called(cart)
	return args.Error(0)
}

func (m *CartRepositoryMock) Update(cart cart.Cart) error {
	args := m.Called(cart)
	return args.Error(0)
}

func (m *CartRepositoryMock) FindOneByCustomerId(customerId uuid.UUID) (*cart.Cart, error) {
	args := m.Called(customerId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cart.Cart), args.Error(1)
}

type AddProductToCartSuite struct {
	suite.Suite
	addProductToCart    usecases.AddProductToCart
	customerGatewayMock CustomerGatewayMock
	cartRepositoryMock  CartRepositoryMock
}

func (a *AddProductToCartSuite) SetupTest() {
	a.customerGatewayMock = CustomerGatewayMock{}
	a.cartRepositoryMock = CartRepositoryMock{}

	a.addProductToCart = usecases.AddProductToCart{
		CustomerGateway: &a.customerGatewayMock,
		CartRepository:  &a.cartRepositoryMock,
	}
}

func (a *AddProductToCartSuite) Test_add_product_to_new_cart_should_succeed() {
	customerId := uuid.New()
	productId := uuid.New()
	a.customerGatewayMock.On("ExistsByCustomerId", mock.Anything).Return(true, nil)
	a.cartRepositoryMock.On("FindOneByCustomerId", mock.Anything).Return(nil, nil)
	a.cartRepositoryMock.On("Create", mock.Anything).Return(nil)
	input := usecases.AddProductToCartInput{
		CustomerId: customerId,
		ProductId:  productId,
		Price:      int64(2440),
		Quantity:   int16(3),
	}

	err := a.addProductToCart.Execute(input)

	a.Equal(nil, err)
	a.cartRepositoryMock.AssertNumberOfCalls(a.T(), "Create", 1)
	a.cartRepositoryMock.AssertNumberOfCalls(a.T(), "Update", 0)
}

func (a *AddProductToCartSuite) Test_add_product_to_existing_cart_should_succeed() {
	customerId := uuid.New()
	productId := uuid.New()
	customerCart := cart.Cart{
		Id:         uuid.New(),
		CustomerId: customerId,
		Items:      []cart.CartItem{},
	}
	a.customerGatewayMock.On("ExistsByCustomerId", mock.Anything).Return(true, nil)
	a.cartRepositoryMock.On("FindOneByCustomerId", mock.Anything).Return(&customerCart, nil)
	a.cartRepositoryMock.On("Update", mock.Anything).Return(nil)
	input := usecases.AddProductToCartInput{
		CustomerId: customerId,
		ProductId:  productId,
		Price:      int64(2440),
		Quantity:   int16(3),
	}

	err := a.addProductToCart.Execute(input)

	a.Equal(nil, err)
	a.cartRepositoryMock.AssertNumberOfCalls(a.T(), "Create", 0)
	a.cartRepositoryMock.AssertNumberOfCalls(a.T(), "Update", 1)
}

func (a *AddProductToCartSuite) Test_add_product_to_cart_with_customer_that_does_not_exist_should_fail() {
	customerId := uuid.New()
	a.customerGatewayMock.On("ExistsByCustomerId", customerId).Return(false, nil)
	input := usecases.AddProductToCartInput{
		CustomerId: customerId,
		ProductId:  uuid.New(),
		Price:      int64(2440),
		Quantity:   int16(3),
	}

	err := a.addProductToCart.Execute(input)

	a.EqualError(err, "customer not found")
}

func TestAddProductToCart(t *testing.T) {
	suite.Run(t, new(AddProductToCartSuite))
}
