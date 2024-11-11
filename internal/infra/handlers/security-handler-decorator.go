package handlers

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gsaaraujo/ecommerce-go/internal/application/gateways"
	webhttp "github.com/gsaaraujo/ecommerce-go/internal/infra/web-http"
	"github.com/labstack/echo/v4"
)

type SecurityHandlerDecorator struct {
	HttpHandler          IHttpHandler
	SecretManagerGateway gateways.ISecretManagerGateway
}

func (a *SecurityHandlerDecorator) Handle(c echo.Context) error {
	authAccessToken, err := a.SecretManagerGateway.Get("AUTH_ACCESS_TOKEN")
	if err != nil {
		return c.JSON(500, webhttp.NewInternalServerError("Something went wrong. Please try again later."))
	}

	authorizationToken := c.Request().Header.Get("Authorization")
	if authorizationToken == "" {
		return c.JSON(401, webhttp.NewUnauthorizedRequest("Authorization token is missing."))
	}

	parts := strings.Split(authorizationToken, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return c.JSON(401, webhttp.NewUnauthorizedRequest("Invalid authorization token format."))
	}

	rawToken := parts[1]

	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(authAccessToken), nil
	})

	if err != nil {
		return c.JSON(401, webhttp.NewUnauthorizedRequest("Authorization token is invalid."))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["customerId"] == nil {
		return c.JSON(401, webhttp.NewForbiddenRequest("You do not have permission to access this resource."))
	}

	c.Set("customerId", claims["customerId"])
	return a.HttpHandler.Handle(c)
}
