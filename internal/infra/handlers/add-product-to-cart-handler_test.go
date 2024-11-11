package handlers_test

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gsaaraujo/ecommerce-go/internal/application/usecases"
	"github.com/gsaaraujo/ecommerce-go/internal/infra"
	"github.com/gsaaraujo/ecommerce-go/internal/infra/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AddToProductMock struct {
	mock.Mock
}

func (a *AddToProductMock) Execute(input usecases.AddProductToCartInput) error {
	args := a.Called(input)
	return args.Error(0)
}

type AddToProductHandlerSuite struct {
	suite.Suite
	addToProductMock    AddToProductMock
	addToProductHandler handlers.AddProductToCartHandler
}

func (a *AddToProductHandlerSuite) SetupTest() {
	a.addToProductMock = AddToProductMock{}
	a.addToProductHandler = handlers.AddProductToCartHandler{
		Validator:        infra.NewValidator(),
		AddProductToCart: &a.addToProductMock,
	}
}

func (a *AddToProductHandlerSuite) Test_request_add_product_to_cart_with_no_errors_should_succeed() {
	e := echo.New()
	a.addToProductMock.On("Execute", mock.Anything).Return(nil)
	request := httptest.NewRequest("POST", "/", strings.NewReader(`
		{
			"productId": "632ef70b-4184-4704-ad7d-8b8f5dd534d9",
			"quantity": 4
		}
	`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)
	context.Set("customerId", "5ad98fc5-6b0f-45fd-a886-d6a15a63c833")

	a.addToProductHandler.Handle(context)

	a.Equal(200, recorder.Code)
	a.JSONEq(`
	{
		"status": "SUCCESS",
		"statusCode": 200,
		"statusText": "OK",
		"data": null
	}
	`, recorder.Body.String())
}

func (a *AddToProductHandlerSuite) Test_request_add_product_to_cart_with_product_not_found_error_should_fail() {
	e := echo.New()
	a.addToProductMock.On("Execute", mock.Anything).Return(errors.New("product not found"))
	request := httptest.NewRequest("POST", "/", strings.NewReader(`
		{
			"productId": "632ef70b-4184-4704-ad7d-8b8f5dd534d9",
			"quantity": 4
		}
	`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)
	context.Set("customerId", "5ad98fc5-6b0f-45fd-a886-d6a15a63c833")

	a.addToProductHandler.Handle(context)

	a.Equal(404, recorder.Code)
	a.JSONEq(`
	{
		"status": "ERROR",
		"statusCode": 404,
		"statusText": "NOT_FOUND",
		"error": "We couldn't find a product with the ID '632ef70b-4184-4704-ad7d-8b8f5dd534d9'. Please check the product ID and try again."
	}
	`, recorder.Body.String())
}

func (a *AddToProductHandlerSuite) Test_request_add_product_to_cart_should_fail_with_invalid_body() {
	a.addToProductMock.On("Execute", mock.Anything).Return(nil)
	bodiesAndErrors := []map[string]string{
		{
			"body":   `abc`,
			"errors": `["content-type must be application/json."]`,
		},
		{
			"body":   `{}`,
			"errors": `["productId is required", "quantity is required"]`,
		},
		{
			"body": `{
				"productId": null,
				"quantity": null
			}`,
			"errors": `["productId is required", "quantity is required"]`,
		},
		{
			"body": `{
				"productId": "",
				"quantity": null
			}`,
			"errors": `["productId must be uuidv4", "quantity is required"]`,
		},
		{
			"body": `{
				"productId": " ",
				"quantity": -3
			}`,
			"errors": `["productId must be uuidv4", "quantity must be greater than or equal to 1"]`,
		},
		{
			"body": `{
				"productId": "abc",
				"quantity": 0
			}`,
			"errors": `["productId must be uuidv4", "quantity must be greater than or equal to 1"]`,
		},
		// {
		// 	"body": `{
		// 		"productId": "3e19ad32-ff8c-4f5c-8919-d2d458502e4c",
		// 		"quantity": 2.5
		// 	}`,
		// 	"errors": `["quantity must be integer"]`,
		// },
		// {
		// 	"body": `{
		// 		"productId": 123,
		// 		"price": "1",
		// 		"quantity": "2"
		// 	}`,
		// 	"errors": `["productId must be string", "price must be integer", "quantity must be integer"]`,
		// },
	}

	for _, inputAndError := range bodiesAndErrors {
		body := inputAndError["body"]
		errorMessage := inputAndError["errors"]

		e := echo.New()
		request := httptest.NewRequest("POST", "/", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		context := e.NewContext(request, recorder)

		a.addToProductHandler.Handle(context)

		a.Equal(400, recorder.Code)
		a.JSONEq(fmt.Sprintf(`
		{
			"status": "ERROR",
			"statusCode": 400,
			"statusText": "BAD_REQUEST",
			"errors": %s
		}
		`, errorMessage), recorder.Body.String())
	}
}

func TestAddProductToCartHandler(t *testing.T) {
	suite.Run(t, new(AddToProductHandlerSuite))
}
