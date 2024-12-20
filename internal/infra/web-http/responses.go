package webhttp

import "github.com/labstack/echo/v4"

type ResponseError struct {
	Status       string `json:"status"`
	StatusCode   uint16 `json:"statusCode"`
	StatusText   string `json:"statusText"`
	ErrorMessage string `json:"error"`
}

type ResponseErrors struct {
	Status        string   `json:"status"`
	StatusCode    uint16   `json:"statusCode"`
	StatusText    string   `json:"statusText"`
	ErrorMessages []string `json:"errors"`
}

type ResponseSuccess struct {
	Status     string      `json:"status"`
	StatusCode uint16      `json:"statusCode"`
	StatusText string      `json:"statusText"`
	Data       interface{} `json:"data"`
}

func NewOk(c echo.Context, data interface{}) error {
	return c.JSON(200, ResponseSuccess{
		Status:     "SUCCESS",
		StatusCode: 200,
		StatusText: "OK",
		Data:       data,
	})
}

func NewBadRequestValidation(c echo.Context, errorMessages []string) error {
	return c.JSON(400, ResponseErrors{
		Status:        "ERROR",
		StatusCode:    400,
		StatusText:    "BAD_REQUEST",
		ErrorMessages: errorMessages,
	})
}

func NewBadRequest(c echo.Context, errorMessage string) error {
	return c.JSON(400, ResponseError{
		Status:       "ERROR",
		StatusCode:   400,
		StatusText:   "BAD_REQUEST",
		ErrorMessage: errorMessage,
	})
}

func NewUnauthorizedRequest(c echo.Context, errorMessage string) error {
	return c.JSON(401, ResponseError{
		Status:       "ERROR",
		StatusCode:   401,
		StatusText:   "UNAUTHORIZED",
		ErrorMessage: errorMessage,
	})
}

func NewForbiddenRequest(c echo.Context, errorMessage string) error {
	return c.JSON(401, ResponseError{
		Status:       "ERROR",
		StatusCode:   401,
		StatusText:   "FORBIDDEN",
		ErrorMessage: errorMessage,
	})
}

func NewNotFound(c echo.Context, errorMessage string) error {
	return c.JSON(404, ResponseError{
		Status:       "ERROR",
		StatusCode:   404,
		StatusText:   "NOT_FOUND",
		ErrorMessage: errorMessage,
	})
}

func NewConflict(c echo.Context, errorMessage string) error {
	return c.JSON(409, ResponseError{
		Status:       "ERROR",
		StatusCode:   409,
		StatusText:   "CONFLICT",
		ErrorMessage: errorMessage,
	})
}

func NewInternalServerError(c echo.Context, errorMessage string) error {
	return c.JSON(500, ResponseError{
		Status:       "ERROR",
		StatusCode:   500,
		StatusText:   "INTERNAL_SERVER_ERROR",
		ErrorMessage: errorMessage,
	})
}
