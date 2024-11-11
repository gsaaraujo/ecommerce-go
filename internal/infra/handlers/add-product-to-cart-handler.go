package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-go/internal/application/usecases"
	"github.com/gsaaraujo/ecommerce-go/internal/infra"
	webhttp "github.com/gsaaraujo/ecommerce-go/internal/infra/web-http"
	"github.com/labstack/echo/v4"
)

type AddProductToCartHandlerInput struct {
	ProductId *string `json:"productId" validate:"required,uuid4"`
	Quantity  *int    `json:"quantity" validate:"required,gte=1"`
}

type AddProductToCartHandler struct {
	Validator        infra.Validator
	AddProductToCart usecases.IAddProductToCart
}

func (a *AddProductToCartHandler) Handle(c echo.Context) error {
	handlerInput := AddProductToCartHandlerInput{}
	if err := c.Bind(&handlerInput); err != nil {
		return c.JSON(400, webhttp.NewBadRequestValidation([]string{"content-type must be application/json."}))
	}

	errorsMessages := a.Validator.Validate(handlerInput)
	if len(errorsMessages) > 0 {
		return c.JSON(400, webhttp.NewBadRequestValidation(errorsMessages))
	}

	productId, err := uuid.Parse(*handlerInput.ProductId)
	if err != nil {
		return c.JSON(500, webhttp.NewInternalServerError("Something went wrong. Please try again later."))
	}

	if c.Get("customerId") == nil {
		return c.JSON(500, webhttp.NewInternalServerError("Something went wrong. Please try again later."))
	}

	customerId, err := uuid.Parse(c.Get("customerId").(string))
	if err != nil {
		return c.JSON(500, webhttp.NewInternalServerError("Something went wrong. Please try again later."))
	}

	err = a.AddProductToCart.Execute(usecases.AddProductToCartInput{
		CustomerId: customerId,
		ProductId:  productId,
		Quantity:   int32(*handlerInput.Quantity),
	})

	if err != nil {
		switch err.Error() {
		case "product not found":
			return c.JSON(404, webhttp.NewNotFound(fmt.Sprintf(`We couldn't find a product with the ID '%s'. Please check the product ID and try again.`,
				*handlerInput.ProductId)))
		}

		return c.JSON(500, webhttp.NewInternalServerError("Something went wrong. Please try again later."))
	}

	return c.JSON(200, webhttp.NewOk(nil))
}
