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

type AddProductToCartMock struct {
	mock.Mock
}

func (a *AddProductToCartMock) Execute(input usecases.AddProductToCartInput) error {
	args := a.Called(input)
	return args.Error(0)
}

type AddProductToCartHandlerSuite struct {
	suite.Suite
	addProductToCartMock    AddProductToCartMock
	addProductToCartHandler handlers.AddProductToCartHandler
}

func (a *AddProductToCartHandlerSuite) SetupTest() {
	a.addProductToCartMock = AddProductToCartMock{}
	a.addProductToCartHandler = handlers.AddProductToCartHandler{
		Validator:        infra.NewValidator(),
		AddProductToCart: &a.addProductToCartMock,
	}
}

func (a *AddProductToCartHandlerSuite) TestAddProductToCartHandler_Handle_OnNoErrors_ReturnsOk() {
	e := echo.New()
	a.addProductToCartMock.On("Execute", mock.Anything).Return(nil)
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

	a.addProductToCartHandler.Handle(context)

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

func (a *AddProductToCartHandlerSuite) TestAddProductToCartHandler_Handle_OnProductNotFound_ReturnsNotFound() {
	e := echo.New()
	a.addProductToCartMock.On("Execute", mock.Anything).Return(errors.New("product not found"))
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

	a.addProductToCartHandler.Handle(context)

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

func (a *AddProductToCartHandlerSuite) TestAddProductToCartHandler_Handle_OnInvalidBody_ReturnsBadRequest() {
	a.addProductToCartMock.On("Execute", mock.Anything).Return(nil)
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

		a.addProductToCartHandler.Handle(context)

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
	suite.Run(t, new(AddProductToCartHandlerSuite))
}
